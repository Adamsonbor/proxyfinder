package router

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"proxyfinder/internal/api"
	"proxyfinder/internal/auth"
	"proxyfinder/internal/domain"
	gormstorage "proxyfinder/internal/storage/v2/gorm-sotrage"

	"github.com/go-chi/chi/v5"
)

var (
	ErrTypeAssertion = fmt.Errorf("Type assertion error")
	ErrInvalidBody   = fmt.Errorf("Invalid body")
	ErrInvalidParams = fmt.Errorf("Invalid params")
)

type FavoritsDTO struct {
	UserId  int64 `json:"user_id"`
	ProxyId int64 `json:"proxy_id"`
}

type FavoritsRouter struct {
	log     *slog.Logger
	storage *gormstorage.Storage
	Router  *chi.Mux
	jwt     auth.JWTService
}

func NewFavoritsRouter(
	log *slog.Logger,
	r *chi.Mux,
	storage *gormstorage.Storage,
	jwt auth.JWTService,
) *FavoritsRouter {
	fr := &FavoritsRouter{
		Router:  r,
		log:     log,
		storage: storage,
		jwt:     jwt,
	}

	r.Route("/favorits", func(r chi.Router) {
		r.Use(jwt.JWTMiddleware)
		r.Route("/{id}", func(r chi.Router) {
			r.Use(idPermissionMiddleware)
			r.Delete("/", fr.DeleteFavorite)
		})
		r.Get("/", fr.GetFavorits)
		r.Post("/", fr.CreateFavorite)
	})

	return fr
}

func (self *FavoritsRouter) GetFavorits(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "FavoritsRouter.GetBy"))

	userId := r.Context().Value("user_id").(int64)

	var favorits []domain.Favorits
	err := self.storage.GetAllBy(&favorits, "user_id", userId)
	if err != nil {
		api.ReturnError(log, w, http.StatusInternalServerError, err)
		return
	}
	log.Info("success", slog.Any("favorits", favorits))

	JSONResponse(w, "success", favorits, err)
}

// /favorits
func (self *FavoritsRouter) CreateFavorite(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "FavoritsRouter.Create"))

	defer r.Body.Close()

	var body domain.Favorits
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Debug("Decode body error")
		api.ReturnError(log, w, http.StatusInternalServerError, err)
		return
	}

	if body.UserId < 1 || body.ProxyId < 1 {
		log.Debug("Invalid body", slog.Int64("userId", body.UserId), slog.Int64("proxyId", body.ProxyId))
		api.ReturnError(log, w, http.StatusInternalServerError, ErrInvalidBody)
		return
	}

	_, err = self.storage.Create(&body)
	if err != nil {
		log.Debug("storage create error", slog.Int64("user_id", body.UserId), slog.Int64("proxy_id", body.ProxyId))
		api.ReturnError(log, w, http.StatusInternalServerError, err)
		return
	}

	JSONResponse(w, "success", nil, err)
}

func (self *FavoritsRouter) DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "FavoritsRouter.Delete"))

	defer r.Body.Close()

	proxyId := r.Context().Value("id").(int64)
	userId := r.Context().Value("user_id").(int64)
	if proxyId < 1 || userId < 1 {
		log.Debug("Invalid body", slog.Int64("userId", userId), slog.Int64("proxyId", proxyId))
		api.ReturnError(log, w, http.StatusInternalServerError, ErrInvalidParams)
		return
	}

	err := self.storage.Delete(&domain.Favorits{}, "proxy_id = ?", proxyId, "user_id = ?", userId)
	if err != nil {
		log.Debug("storage delete error", slog.Int64("user_id", userId), slog.Int64("proxy_id", proxyId))
		api.ReturnError(log, w, http.StatusInternalServerError, err)
		return
	}

	JSONResponse(w, "success", nil, err)
}
