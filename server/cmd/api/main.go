package main

import (
	"fmt"
	"net/http"
	googleapi "proxyfinder/internal/api/chi-router/auth/google"
	rabbitApi "proxyfinder/internal/api/chi-router/rabbit"
	router "proxyfinder/internal/api/chi-router/v1/gorm"
	routerv2 "proxyfinder/internal/api/chi-router/v2/gorm-sqlx"
	googleauth "proxyfinder/internal/auth/google"
	"proxyfinder/internal/broker/rabbit"
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
	fmt.Println(cfg)

	// INIT logger
	log := logger.New(cfg.Env)
	log.Info("Initializing with env: " + cfg.Env)

	// INIT gorm sqlite
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
	sqlxdb.SetMaxIdleConns(10)
	sqlxdb.SetMaxOpenConns(100)
	sqlxdb.SetConnMaxLifetime(5)

	// INIT sqlx storage
	sqlxStorage := sqlxstorage.New(sqlxdb)

	// INIT rabbitmq
	rabbitService := rabbit.NewRabbit(cfg, "mail")

	// INIT google auth
	googleAuth := googleauth.NewGoogleAuth(&cfg.GoogleAuth)

	// INIT router
	mux := chi.NewMux()

	routerv1 := router.New(log, storage)
	routerv2 := routerv2.New(log, storagev2, sqlxStorage)
	routerRabbit := rabbitApi.New(log, rabbitService, cfg)
	routerGoogleAuth := googleapi.NewRouter(log, cfg, googleAuth)

	mux.Mount("/api/v1", routerv1.Router)
	mux.Mount("/api/v2", routerv2.Router)
	mux.Mount("/rabbit", routerRabbit.Router)
	mux.Mount("/auth/google", routerGoogleAuth.Router)

	for _, route := range mux.Routes() {
		log.Debug(fmt.Sprintf("%s", route.Pattern))
	}

	http.ListenAndServe(":8080", mux)
}
