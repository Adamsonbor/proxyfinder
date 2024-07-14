package main

import (
	"context"
	"proxyfinder/internal/collector/geonode"
	"proxyfinder/internal/config"
	"proxyfinder/internal/logger"
	"proxyfinder/internal/storage/sqlite-storage"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config.MustLoadConfig()

	db, err := sqlx.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.Env)

	countryStorage := sqlite.NewCountry(db)
	proxyStorage := sqlite.NewProxy(db)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Collector.Timeout)
	defer cancel()

	_, err = geonode.New(log, proxyStorage, countryStorage).Collect(ctx)
	if err != nil {
		panic(err)
	}

	// fmt.Println(proxies)
}
