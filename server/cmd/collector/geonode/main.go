package main

import (
	"context"
	"fmt"
	"log/slog"
	"proxyfinder/internal/collector/geonode"
	"proxyfinder/internal/config"
	"proxyfinder/internal/logger"
	"proxyfinder/internal/storage/sqlx-storage"

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

	countryStorage := sqlxstorage.NewCountry(db)
	proxyStorage := sqlxstorage.NewProxy(db)

	collector := geonode.New(log, proxyStorage, countryStorage)
	pager := collector.NewPageScheduler()

	page := pager()
	// fmt.Println(page)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Collector.Timeout)
	defer cancel()

	for i := range 10 {
		filename := fmt.Sprintf("storage/init/proxies_%d.json", i)
		log := log.With(slog.String("filename", filename), slog.String("page", page))
		_, err = collector.Collect(ctx, page, filename)
		log.Info("collector.Collect done")
		if err != nil {
			log.Error("collector.Collect failed", slog.String("err", err.Error()))
			return
		}
	}
}
