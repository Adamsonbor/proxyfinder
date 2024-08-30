package userapi

import (
	"context"
	"log/slog"
	"net/http"
	"proxyfinder/internal/config"
	apiv1 "proxyfinder/internal/service/api"
	serviceapiv1 "proxyfinder/internal/service/api"
	jwtservice "proxyfinder/internal/service/api/v1/auth/jwt"
	chiapi "proxyfinder/internal/transport/api/chi"
	chiapiv1 "proxyfinder/internal/transport/api/chi/v1"

	"github.com/go-chi/chi/v5"
)

type UserController struct {
	log     *slog.Logger
	Router  *chi.Mux
	service apiv1.UserService
	jwt     serviceapiv1.JWTService
	cfg     config.Config
}

func New(
	log *slog.Logger,
	service apiv1.UserService,
	jwt serviceapiv1.JWTService,
	cfg config.Config,
) *UserController {
	r := chi.NewRouter()
	router := &UserController{
		log:     log,
		Router:  r,
		service: service,
		jwt:     jwt,
		cfg:     cfg,
	}

	r.Use(chiapiv1.FilterMiddleware)
	r.Use(chiapiv1.SortMiddleware)
	r.Use(jwt.JWTMiddleware)
	r.Get("/", router.UserInfo)

	return router
}

// UserInfo
// @Summary Get user info
// @Description Get user info
// @Tags user
// @Param Authorization header string true "Authorization"
// @Produce json
// @Success 200 {object} chiapi.Response
// @Failure 400 {object} chiapi.Response
// @Failure 500 {object} chiapi.Response
// @Router /api/v1/user [get]
func (self *UserController) UserInfo(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "user.UserInfo"))

	id := r.Context().Value(jwtservice.JwtUserCtxKey).(int64)

	ctx, cancel := context.WithTimeout(r.Context(), self.cfg.Database.Timeout)
	defer cancel()

	user, err := self.service.UserInfo(ctx, id)
	if err != nil {
		log.Debug("failed to get user info", slog.Any("error", err))
		chiapi.JSONresponse(w, http.StatusBadRequest, nil, err)
		return
	}

	log.Debug("user", slog.Any("user", user))

	chiapi.JSONresponse(w, http.StatusOK, user, nil)
}
