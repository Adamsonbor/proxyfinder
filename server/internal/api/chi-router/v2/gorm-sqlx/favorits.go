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

type FavoritsRouter struct {
	log     *slog.Logger
	storage *gormstorage.Storage
	Router  *chi.Mux
	jwt     auth.JWTService
}

type CreateFavoriteRequest struct {
	UserId  int64
	ProxyId int64
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
			r.Delete("/", fr.DeleteFavorite)
		})
		r.Get("/", fr.GetFavorits)
	})
	r.Post("/favorits", fr.CreateFavorite)

	return fr
}

func (self *FavoritsRouter) GetFavorits(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "FavoritsRouter.GetBy"))

	userId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		log.Error("get", slog.String("err", "user_id not found"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var favorits []domain.Favorits
	err := self.storage.GetAllBy(&favorits, "userId", userId)
	if err != nil {
		log.Error("get", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info("success", slog.Any("favorits", favorits))

	JSONResponse(w, "success", favorits, err)
}

// /favorits
func (self *FavoritsRouter) CreateFavorite(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "FavoritsRouter.Create"))

	var req CreateFavoriteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error("decode", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	JSONResponse(w, "success", nil, err)
}

func (self *FavoritsRouter) DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "FavoritsRouter.Delete"))

	id := r.Context().Value("id").(int64)

	err := self.storage.Delete(domain.Favorits{Id: id})
	if err != nil {
		log.Error("delete", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSONResponse(w, "success", nil, err)
}
