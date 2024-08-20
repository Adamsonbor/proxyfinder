package googleapi

import (
	"context"
	"log/slog"
	"net/http"

	// dependencies
	"proxyfinder/internal/config"
	"proxyfinder/internal/domain"
	serviceapiv1 "proxyfinder/internal/service/api"
	transportapi "proxyfinder/internal/transport/api"

	"time"

	"github.com/go-chi/chi/v5"
)

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type Router struct {
	log     *slog.Logger
	Router  *chi.Mux
	service serviceapiv1.GoogleAuthService
	cfg     config.Config
}

func New(
	log *slog.Logger,
	service serviceapiv1.GoogleAuthService,
	cfg config.Config,
) *Router {
	r := chi.NewRouter()

	router := &Router{
		log:     log,
		service: service,
		Router:  r,
		cfg:     cfg,
	}

	// redirect to google login and redirect to callback with code (query param)
	r.Get("/login", router.Login)

	// get token from google using code and create new user using google user info
	// create new jwt (access_token, refresh_token) and set it in cookie
	// and redirect to frontend
	r.Get("/callback", router.Callback)

	// get refresh token from query params and create new jwt pair (access_token, refresh_token)
	// and set it in cookie
	r.Get("/refresh", router.Refresh)

	return router
}

// Get refresh token from query params and create new jwt pair (access_token, refresh_token)
// and return it in json
// @Summary update refresh token
// @Description update refresh token
// @Tags auth
// @Param refresh_token query string true "refresh token"
// @Success 200
// @Router /refresh [get]
func (self *Router) Refresh(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "GoogleAuth.Refresh"))

	// validate token
	refreshToken := r.URL.Query().Get("refresh_token")

	// update refresh token
	ctx, cancel := context.WithTimeout(r.Context(), self.cfg.GoogleAuth.Timeout)
	defer cancel()

	res, err := self.service.UpdateRefreshToken(ctx, refreshToken)
	if err != nil {
		log.Debug("update refresh token error", slog.Any("error", err))
		transportapi.JSONresponse(w, http.StatusInternalServerError, nil, err)
		return
	}

	transportapi.JSONresponse(w, http.StatusOK, res, nil)
}

// @Summary redirect to google login
// @Description redirect to google login
// @Tags auth
// @Success 200
// @Router /login [get]
func (self *Router) Login(w http.ResponseWriter, r *http.Request) {
	url := self.service.Login("state")
	http.Redirect(w, r, url, http.StatusFound)
}

// Get access token from google, save new user in database,
// create new jwt token and set it in cookie
// and redirect to frontend
func (self *Router) Callback(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "GoogleAuth.Callback"))
	log.Debug("request", slog.Any("request", r.URL.Query()))

	// get access token
	code := r.FormValue("code")

	// Get or create new user and
	ctx, cancel := context.WithTimeout(r.Context(), self.cfg.GoogleAuth.Timeout)
	defer cancel()
	user, err := self.service.Callback(ctx, code)
	if err != nil {
		log.Debug("callback error", slog.Any("error", err))
		transportapi.JSONresponse(w, http.StatusInternalServerError, nil, err)
		return
	}
	log.Debug("user", slog.Any("user", user))

	// generate new tokens and set coockies
	err = self.GenerateAndSetCoockies(w, r, user)
	if err != nil {
		log.Debug("generate and set coockies error", slog.Any("error", err))
		transportapi.JSONresponse(w, http.StatusInternalServerError, nil, err)
		return
	}

	// redirect to frontend
	http.Redirect(w, r, self.cfg.GoogleAuth.RedirectTo, http.StatusFound)
}

func (self *Router) GenerateAndSetCoockies(w http.ResponseWriter, r *http.Request, user domain.User) error {
	accessToken, refreshToken, err := self.service.GenerateTokens(user)
	if err != nil {
		return err
	}

	// set cookies
	http.SetCookie(w, &http.Cookie{
		Name:    "access_token",
		Value:   accessToken,
		Expires: time.Now().Add(time.Duration(self.cfg.JWT.AccessTokenTTL.Seconds())),
		Path:    "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "refresh_token",
		Value:   refreshToken,
		Expires: time.Now().Add(time.Duration(self.cfg.JWT.RefreshTokenTTL.Seconds())),
		Path:    "/",
	})

	return nil
}

