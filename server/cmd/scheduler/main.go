package main

import (
	"fmt"
	httpchecker "proxyfinder/internal/service/checker/http-checker"
	defaultscheduler "proxyfinder/internal/service/scheduler/default-scheduler"
	proxystorage "proxyfinder/internal/storage/sqlx/proxy"
	"proxyfinder/pkg/logger"

	"proxyfinder/internal/config"

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
	proxyStorage := proxystorage.New(db)
	log.Info("ProxyStorage connected")

	// INIT scheduler
	scheduler := defaultscheduler.New(cfg, log, proxyStorage, checker)
	log.Info("Scheduler connected")

	// START
	scheduler.Run()
}
