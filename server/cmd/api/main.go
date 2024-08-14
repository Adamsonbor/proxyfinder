package main

import (
	"log/slog"
	"net/http"
	googleapi "proxyfinder/internal/api/chi-router/auth/google"
	// rabbitApi "proxyfinder/internal/api/chi-router/rabbit"
	router "proxyfinder/internal/api/chi-router/v1/gorm"
	routerv2 "proxyfinder/internal/api/chi-router/v2/gorm-sqlx"
	jwtservice "proxyfinder/internal/auth/jwt"
	// "proxyfinder/internal/broker/rabbit"
	"proxyfinder/internal/config"
	"proxyfinder/internal/logger"
	gormstoragev1 "proxyfinder/internal/storage/gorm-storage"
	gormstoragev2 "proxyfinder/internal/storage/v2/gorm-sotrage"
	sqlxstorage "proxyfinder/internal/storage/v2/sqlx-storage"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	// INIT config
	cfg := config.MustLoadConfig()

	// INIT logger
	log := logger.New(cfg.Env)
	log.Info("Initializing with env: " + cfg.Env)
	db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// INIT storage
	storage := gormstoragev1.New(db)
	storagev2 := gormstoragev2.New(db)

	// INIT sqlxdb
	sqlxdb, err := sqlx.Connect("sqlite3", cfg.Database.Path)
	if err != nil {
		panic(err)
	}

	// INIT sqlx storage
	sqlxStorage := sqlxstorage.New(sqlxdb)

	// // INIT rabbitmq
	// rabbitService := rabbit.NewRabbit(cfg, "mail")

	// INIT jwt
	jwt := jwtservice.NewJWTService(log, cfg, sqlxStorage.UserStorage)

	// INIT router
	mux := chi.NewMux()

	// INIT routers
	routerv1 := router.New(log, storage)
	routerv2 := routerv2.New(log, storagev2, sqlxStorage, jwt)
	// routerRabbit := rabbitApi.New(log, rabbitService, cfg)
	routerGoogle := googleapi.NewRouter(log, cfg, sqlxStorage.UserStorage)

	// register routes
	mux.Mount("/api/v1", routerv1.Router)
	mux.Mount("/api/v2", routerv2.Router)
	// mux.Mount("/rabbit", routerRabbit.Router)
	mux.Mount("/auth/google", routerGoogle.Router)

	// print routes
	chi.Walk(mux, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Info(route, slog.String("method", method))
		return nil
	})

	// run server
	http.ListenAndServe(":8080", mux)
}
