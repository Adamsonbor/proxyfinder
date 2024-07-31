package main

import (
	"fmt"
	httpchecker "proxyfinder/internal/checker/http-checker"
	defaultScheduler "proxyfinder/internal/scheduler/default-scheduler"
	"proxyfinder/internal/storage/sqlx-storage"

	"proxyfinder/internal/config"
	"proxyfinder/internal/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// INIT config
	cfg := config.MustLoadConfig()

	// INIT logger
	log := logger.New(cfg.Env)
	log.Info(fmt.Sprintf("Initializing with env: %s", cfg.Env))

	// INIT checker
	checker := httpchecker.New(log)

	log.Info("Database path: " + cfg.Database.Path)
	// INIT database
	db, err := sqlx.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		panic(err)
	}
	log.Info("Database connected")

	// INIT ProxyStorage
	proxyStorage := sqlxstorage.NewProxy(db)
	log.Info("ProxyStorage connected")

	// INIT scheduler
	scheduler := defaultScheduler.NewScheduler(cfg, log, proxyStorage, checker)
	log.Info("Scheduler connected")

	// START
	scheduler.Run()
}
