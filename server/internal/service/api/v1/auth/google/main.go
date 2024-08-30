package googleservice

import (
	"context"
	"log/slog"
	"proxyfinder/internal/config"
	"proxyfinder/internal/domain"
	serviceapiv1 "proxyfinder/internal/service/api"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

type GoogleAuthService struct {
	log         *slog.Logger
	userService serviceapiv1.UserService
	jwt         serviceapiv1.JWTService
	cfg         *config.Config
	gCfg        *oauth2.Config
}

func New(
	log *slog.Logger,
	userService serviceapiv1.UserService,
	jwt serviceapiv1.JWTService,
	cfg *config.Config,
) serviceapiv1.GoogleAuthService {
	return &GoogleAuthService{
		log:         log,
		userService: userService,
		jwt:         jwt,
		cfg:         cfg,
		gCfg: &oauth2.Config{
			ClientID:     cfg.GoogleAuth.ClientId,
			ClientSecret: cfg.GoogleAuth.ClientSecret,
			RedirectURL:  cfg.GoogleAuth.RedirectUrl,
			Scopes:       cfg.GoogleAuth.Scope,
			Endpoint:     google.Endpoint,
		},
	}
}

func (self *GoogleAuthService) Login(state string) string {
	return self.gCfg.AuthCodeURL(state)
}

func (self *GoogleAuthService) UpdateRefreshToken(
	ctx context.Context,
	token string,
) (*serviceapiv1.JWTokens, error) {
	log := self.log.With(slog.String("op", "GoogleAuthService.UpdateRefreshToken"))

	err := self.jwt.ValidateToken(token)
	if err != nil {
		log.Debug("validate token error", slog.Any("error", err))
		return nil, err
	}
	log.Debug("token is valid")

	// get user by refresh token
	user, err := self.userService.GetBy(ctx, "refresh_token", token)
	if err != nil {
		log.Debug("get user by refresh token error", slog.Any("error", err))
		return nil, err
	}
	log.Debug("user", slog.Any("user", user))

	// generate new tokens
	access, refresh, err := self.GenerateTokens(user)
	if err != nil {
		log.Debug("generate tokens error", slog.Any("error", err))
		return nil, err
	}
	log.Debug("tokens", slog.Any("access", access), slog.Any("refresh", refresh))

	// sign tokens
	accessString, err := access.SignedString([]byte(self.cfg.JWT.Secret))
	if err != nil {
		log.Debug("sign access token error", slog.Any("error", err))
		return nil, err
	}
	refreshString, err := refresh.SignedString([]byte(self.cfg.JWT.Secret))
	if err != nil {
		log.Debug("sign refresh token error", slog.Any("error", err))
		return nil, err
	}
	log.Debug("tokens signed")

	// create new session
	err = self.userService.NewSession(ctx, user.Id, refreshString, refresh.Claims.(jwt.MapClaims)["exp"].(int64))
	if err != nil {
		log.Debug("new session error", slog.Any("error", err))
		return nil, err
	}

	log.Debug("session created")

	// return new tokens in json
	res := &serviceapiv1.JWTokens{
		AccessToken:  accessString,
		RefreshToken: refreshString,
		ExpiresIn:    access.Claims.(jwt.MapClaims)["exp"].(int64),
		ExpiresInRef: refresh.Claims.(jwt.MapClaims)["exp"].(int64),
	}
	return res, nil
}

func (self *GoogleAuthService) Callback(
	ctx context.Context,
	code string,
) (*serviceapiv1.JWTokens, error) {
	log := self.log.With(slog.String("op", "GoogleAuthService.Callback"))

	token, err := self.gCfg.Exchange(ctx, code)
	if err != nil {
		log.Debug("exchange error", slog.Any("error", err))
		return nil, err
	}
	log.Debug("token", slog.Any("token", token))

	// get user info from google api
	userInfo, err := self.UserInfo(token)
	if err != nil {
		log.Debug("get user info error", slog.Any("error", err))
		return nil, err
	}
	log.Debug("user info", slog.Any("user", userInfo))

	// Get or create user
	user, err := self.userService.GetBy(ctx, "email", userInfo.Email)
	log.Debug("get or create user", slog.Any("user", user), slog.Any("error", err))
	if err != nil && err.Error() == serviceapiv1.ErrRecordNotFound {
		id, err := self.userService.Save(ctx, *userInfo)
		if err != nil {
			log.Debug("create user error", slog.Any("error", err))
			return nil, err
		}
		userInfo.Id = id
		user = *userInfo
	} else if err != nil {
		log.Debug("get user error", slog.Any("error", err))
		return nil, err
	}

	// generate new tokens
	access, refresh, err := self.GenerateTokens(user)
	if err != nil {
		log.Debug("generate tokens error", slog.Any("error", err))
		return nil, err
	}
	log.Debug("tokens", slog.Any("access", access), slog.Any("refresh", refresh))

	// sign tokens
	accessString, refreshString, err := self.SignTokens(access, refresh)
	if err != nil {
		log.Debug("sign tokens error", slog.Any("error", err))
		return nil, err
	}
	log.Debug("tokens signed")

	// create new session
	err = self.userService.NewSession(ctx, user.Id, refreshString, refresh.Claims.(jwt.MapClaims)["exp"].(int64))
	if err != nil {
		log.Debug("new session error", slog.Any("error", err))
		return nil, err
	}

	return &serviceapiv1.JWTokens{
		AccessToken:  accessString,
		RefreshToken: refreshString,
		ExpiresIn:    access.Claims.(jwt.MapClaims)["exp"].(int64),
		ExpiresInRef: refresh.Claims.(jwt.MapClaims)["exp"].(int64),
	}, nil
}

func (self *GoogleAuthService) GenerateTokens(user domain.User) (*jwt.Token, *jwt.Token, error) {
	accessToken, err := self.jwt.GenerateAccessToken(user.Id)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := self.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (self *GoogleAuthService) SignTokens(access *jwt.Token, refresh *jwt.Token) (string, string, error) {
	accessString, err := access.SignedString([]byte(self.cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	refreshString, err := refresh.SignedString([]byte(self.cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, nil
}

// Get all user info from google api using access token
func (self *GoogleAuthService) UserInfo(token *oauth2.Token) (*domain.User, error) {
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
