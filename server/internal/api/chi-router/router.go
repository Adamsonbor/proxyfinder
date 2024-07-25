package chirouter

import (
	"log/slog"
	router "proxyfinder/internal/api/chi-router/v1/gorm"
	routerv2 "proxyfinder/internal/api/chi-router/v2/gorm-sqlx"
	gormstorage "proxyfinder/internal/storage/gorm-storage"
	sqlxstorage "proxyfinder/internal/storage/v2/sqlx-storage"

	"github.com/go-chi/chi/v5"
)

func New(
	log *slog.Logger,
	storage *gormstorage.Storage,
	sqlxStorage *sqlxstorage.ProxyStorage,
) *chi.Mux {

	routerv1 := router.New(log, storage)
	routerv2 := routerv2.New(log, storage, sqlxStorage)

	mux := chi.NewMux()
	mux.Mount("/api/v1", routerv1.Router)
	mux.Mount("/api/v2", routerv2.Router)

	return mux
}
