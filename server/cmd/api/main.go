package main

import (
	"fmt"
	"net/http"
	chirouter "proxyfinder/internal/api/chi-router"
	"proxyfinder/internal/config"
	"proxyfinder/internal/logger"
	gormstorage "proxyfinder/internal/storage/gorm-storage"
	sqlxstorage "proxyfinder/internal/storage/v2/sqlx-storage"

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
	storage := gormstorage.New(db)

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

	// INIT router
	mux := chirouter.New(log, storage, sqlxStorage)

	http.ListenAndServe(":8080", mux)
}
