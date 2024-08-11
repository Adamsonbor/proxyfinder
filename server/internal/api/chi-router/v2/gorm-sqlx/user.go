package router

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"proxyfinder/internal/auth"
	"proxyfinder/internal/domain"
	gormstorage "proxyfinder/internal/storage/v2/gorm-sotrage"

	"github.com/go-chi/chi/v5"
)

type UserRouter struct {
	log     *slog.Logger
	Router  *chi.Mux
	storage *gormstorage.Storage
	jwt     auth.JWTService
}

func NewUserRouter(
	log *slog.Logger,
	r *chi.Mux,
	storage *gormstorage.Storage,
	jwt auth.JWTService,
) *UserRouter {
	ur := &UserRouter{
		log:     log,
		Router:  r,
		storage: storage,
		jwt:     jwt,
	}

	r.Route("/user", func(r chi.Router) {
		r.Use(ur.jwt.JWTMiddleware)
		r.Get("/", ur.GetUser)
		r.Put("/", ur.UpdateUser)
		r.Delete("/", ur.DeleteUser)
	})

	return ur
}

func (self *UserRouter) GetUser(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "UserRouter.GetUser"))

	userId := r.Context().Value("user_id").(int64)

	var user domain.User
	err := self.storage.Get(&user, userId)
	if err != nil {
		log.Error("get", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info("success", slog.Any("user", user))

	JSONResponse(w, "success", user, err)
}

func (self *UserRouter) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inst, err := self.storage.Create(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSONResponse(w, "success", inst, err)
}

func (self *UserRouter) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.Id = r.Context().Value("user_id").(int64)

	inst, err := self.storage.Update(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSONResponse(w, "success", inst, err)
}

func (self *UserRouter) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(int64)

	err := self.storage.Delete(&domain.User{Id: userId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSONResponse(w, "success", nil, err)
}
