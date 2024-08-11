package googleapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"proxyfinder/internal/api"
	jwtservice "proxyfinder/internal/auth/jwt"
	"proxyfinder/internal/config"
	"proxyfinder/internal/domain"
	gormstorage "proxyfinder/internal/storage/v2/gorm-sotrage"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
	"gorm.io/gorm"
)

type Router struct {
	log         *slog.Logger
	cfg         *config.Config
	gCfg        *oauth2.Config
	userStorage *gormstorage.Storage
	jwt         *jwtservice.JWTService
	Router      *chi.Mux
}

func NewRouter(
	log *slog.Logger,
	cfg *config.Config,
	userStorage *gormstorage.Storage,
) *Router {
	r := chi.NewRouter()

	router := &Router{
		log: log,
		cfg: cfg,
		gCfg: &oauth2.Config{
			ClientID:     cfg.GoogleAuth.ClientId,
			ClientSecret: cfg.GoogleAuth.ClientSecret,
			RedirectURL:  cfg.GoogleAuth.RedirectUrl,
			Scopes: []string{
				people.UserinfoProfileScope,
				people.UserinfoEmailScope,
			},
			Endpoint: google.Endpoint,
		},
		userStorage: userStorage,
		jwt:         jwtservice.NewJWTService(log, cfg, userStorage),
		Router:      r,
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		ExposedHeaders: []string{"Content-Range"},
		MaxAge:         300, // Maximum value not ignored by any of major browsers
	}))

	// redirect to google login and redirect to callback with code (query param)
	r.Get("/login", router.Login)

	// get token from google using code and create new user using google user info
	// create new jwt (access_token, refresh_token) and set it in cookie
	// and redirect to frontend
	r.Get("/callback", router.Callback)

	return router
}

func (self *Router) Login(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, self.gCfg.AuthCodeURL("state"), http.StatusFound)
}

func (self *Router) Callback(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "GoogleAuth.Callback"))
	log.Debug("request", slog.Any("request", r.URL.Query()))

	// get access token
	code := r.FormValue("code")
	token, err := self.gCfg.Exchange(r.Context(), code)
	if err != nil {
		ReturnError(log, w, http.StatusUnauthorized, err)
		return
	}
	log.Debug("token", slog.Any("token", token))

	// get user info from google api
	user, err := self.UserInfo(token)
	if err != nil {
		ReturnError(log, w, http.StatusUnauthorized, err)
		return
	}
	log.Debug("user info", slog.Any("user", user))

	// check if user exists
	// if not create new
	err = self.userStorage.GetBy(user, "email", user.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		inst, err := self.userStorage.Create(user)
		log.Debug("create user", slog.Any("user", inst))
		if err != nil {
			ReturnError(log, w, http.StatusInternalServerError, err)
			return
		}

		newUser, ok := inst.(*domain.User)
		if !ok {
			ReturnError(log, w, http.StatusInternalServerError, err)
			return
		}

		user = newUser
	}

	log.Debug("success", slog.Any("user", user))
	// generate new tokens and set it in cookies
	err = self.GenerateAndSetCoockies(w, r, user)
	if err != nil {
		ReturnError(log, w, http.StatusInternalServerError, err)
		return
	}

	// redirect to frontend
	http.Redirect(w, r, self.cfg.GoogleAuth.RedirectTo, http.StatusFound)
}

func (self *Router) GenerateAndSetCoockies(w http.ResponseWriter, r *http.Request, user *domain.User) error {
	accessToken, refreshToken, err := self.GenerateTokens(user)
	if err != nil {
		return err
	}
	self.setCookies(w, accessToken, refreshToken)

	return nil
}

func (self *Router) GenerateTokens(user *domain.User) (string, string, error) {
	accessToken, err := self.jwt.GenerateAccessToken(user.Id)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := self.jwt.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (self *Router) setCookies(w http.ResponseWriter, accessToken string, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:    "access_token",
		Value:   accessToken,
		Expires: time.Now().Add(self.cfg.JWT.AccessTokenTTL),
		MaxAge:  int(self.cfg.JWT.AccessTokenTTL.Seconds()),
		Path:    "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "refresh_token",
		Value:   refreshToken,
		Expires: time.Now().Add(self.cfg.JWT.RefreshTokenTTL),
		MaxAge:  int(self.cfg.JWT.RefreshTokenTTL.Seconds()),
		Path:    "/",
	})
}

// Get all user info from google api using access token
func (self *Router) UserInfo(token *oauth2.Token) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), self.cfg.JWT.Timeout)
	defer cancel()

	client := self.gCfg.Client(ctx, token)

	svc, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	userInfo, err := svc.People.Get("people/me").PersonFields("names,emailAddresses,photos,phoneNumbers,birthdays").Do()
	if err != nil {
		return nil, err
	}

	return PeopleToUser(userInfo), nil
}

// convert google user info to domain user
func PeopleToUser(userInfo *people.Person) *domain.User {
	user := &domain.User{}
	if len(userInfo.EmailAddresses) > 0 {
		user.Email = userInfo.EmailAddresses[0].Value
	}

	if len(userInfo.Names) > 0 {
		user.Name = userInfo.Names[0].DisplayName
	}

	if len(userInfo.Photos) > 0 {
		user.PhotoUrl = userInfo.Photos[0].Url
	}

	if len(userInfo.PhoneNumbers) > 0 {
		user.Phone = userInfo.PhoneNumbers[0].Value
	}

	if len(userInfo.Birthdays) > 0 {
		timeDate := time.Date(
			int(userInfo.Birthdays[0].Date.Year),
			time.Month(userInfo.Birthdays[0].Date.Month),
			int(userInfo.Birthdays[0].Date.Day),
			0, 0, 0, 0, time.UTC,
		)
		user.DateOfBirth = timeDate
	}

	return user
}

func ReturnError(log *slog.Logger, w http.ResponseWriter, statusCode int, err error) {
	log.Error("error", slog.Any("error", err))
	w.WriteHeader(statusCode)
	api.ReturnResponse(w, "error", nil, err.Error())
}
