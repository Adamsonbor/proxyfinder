package googleapi

import (
	"context"
	"log/slog"
	"net/http"

	// dependencies
	"proxyfinder/internal/config"
	serviceapiv1 "proxyfinder/internal/service/api"
	chiapi "proxyfinder/internal/transport/api/chi"

	"time"

	"github.com/go-chi/chi/v5"
)

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
// @Router /auth/google/refresh [get]
func (self *Router) Refresh(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "GoogleAuth.Refresh"))

	// validate token
	refreshToken := r.URL.Query().Get("refresh_token")

	log.Debug("refresh token", slog.String("token", refreshToken))

	// update refresh token
	ctx, cancel := context.WithTimeout(r.Context(), self.cfg.GoogleAuth.Timeout)
	defer cancel()

	res, err := self.service.UpdateRefreshToken(ctx, refreshToken)
	if err != nil {
		log.Debug("update refresh token error", slog.Any("error", err))
		chiapi.JSONresponse(w, http.StatusInternalServerError, nil, err)
		return
	}

	log.Debug("tokens", slog.Any("token", res))

	chiapi.JSONresponse(w, http.StatusOK, res, nil)
}

// @Summary redirect to google login
// @Description redirect to google login
// @Tags auth
// @Success 200
// @Router /auth/google/login [get]
func (self *Router) Login(w http.ResponseWriter, r *http.Request) {
	url := self.service.Login("state")
	http.Redirect(w, r, url, http.StatusFound)
}

// Get access token from google, save new user in database,
// create new jwt token and set it in cookie
// and redirect to frontend
func (self *Router) Callback(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "GoogleAuth.Callback"))

	if r.URL.Query().Get("error") != "" {
		log.Debug("google access denied", slog.Any("error", r.URL.Query().Get("error")))
		http.Redirect(w, r, self.cfg.GoogleAuth.RedirectTo, http.StatusFound)
		return
	}

	log.Debug("request", slog.Any("request", r.URL.Query()))


	// get access token from google
	code := r.FormValue("code")
	log.Debug("code", slog.String("code", code))

	// service return serviceapiv1.JWTokens
	ctx, cancel := context.WithTimeout(r.Context(), self.cfg.GoogleAuth.Timeout)
	defer cancel()
	tokens, err := self.service.Callback(ctx, code)
	if err != nil {
		log.Debug("callback error", slog.Any("error", err))
		http.Redirect(w, r, self.cfg.GoogleAuth.RedirectTo, http.StatusFound)
		return
	}

	// set cookies
	err = self.SetCookies(w, tokens)
	if err != nil {
		log.Debug("set cookies error", slog.Any("error", err))
		http.Redirect(w, r, self.cfg.GoogleAuth.RedirectTo, http.StatusFound)
		return
	}

	// redirect to frontend with cookies
	log.Debug("redirect to frontend")
	http.Redirect(w, r, self.cfg.GoogleAuth.RedirectTo, http.StatusFound)
}

func (self *Router) SetCookies(w http.ResponseWriter, tokens *serviceapiv1.JWTokens) error {
	// set cookies
	http.SetCookie(w, &http.Cookie{
		Name:    "access_token",
		Value:   tokens.AccessToken,
		Expires: time.Unix(tokens.ExpiresIn, 0).UTC(),
		Path:    "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "refresh_token",
		Value:   tokens.RefreshToken,
		Expires: time.Unix(tokens.ExpiresIn, 0).UTC(),
		Path:    "/",
	})

	return nil
}
