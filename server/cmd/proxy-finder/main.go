package main

import (
	"fmt"
	httpchecker "proxyfinder/internal/checker/http-checker"
	defaultScheduler "proxyfinder/internal/scheduler/default-scheduler"
	"proxyfinder/internal/storage/sqlite-storage"

	"proxyfinder/internal/config"
	"proxyfinder/internal/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// INIT database
	db, err := sqlx.Open("sqlite3", "./storage/local.db")
	if err != nil {
		panic(err)
	}

	// INIT config
	cfg := config.MustLoadConfig()

	// INIT logger
	log := logger.New(cfg.Env)
	log.Info(fmt.Sprintf("Initializing with env: %s", cfg.Env))

	// INIT checker
	checker := httpchecker.New(log)

	// INIT ProxyStorage
	proxyStorage := sqlite.NewProxy(db)

	// INIT scheduler
	scheduler := defaultScheduler.NewScheduler(cfg, log, proxyStorage, checker)

	// START
	scheduler.Run()
}
