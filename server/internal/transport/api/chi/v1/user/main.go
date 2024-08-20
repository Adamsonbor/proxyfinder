package userapi

import (
	"log/slog"
	"net/http"
	apiv1 "proxyfinder/internal/service/api"
	chiapiv1 "proxyfinder/internal/transport/api/chi/v1"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserController struct {
	log     *slog.Logger
	Router  *chi.Mux
	service apiv1.UserService
}

func New(log *slog.Logger, service apiv1.UserService) *UserController {
	r := chi.NewRouter()
	router := &UserController{
		log:     log,
		Router:  r,
		service: service,
	}

	r.Use(chiapiv1.FilterMiddleware)
	r.Get("/user/{id}", router.UserInfo)

	return router
}

func (self *UserController) UserInfo(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "user.UserInfo"))
	idStr := chi.URLParam(r, "id")
	id64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := self.service.UserInfo(id)
	if err != nil {
	}
}
