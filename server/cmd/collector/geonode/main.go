package main

import (
	"context"
	"fmt"
	"log/slog"
	"proxyfinder/internal/collector/geonode"
	"proxyfinder/internal/config"
	"proxyfinder/internal/logger"

	"github.com/google/uuid"
)

func main() {
	cfg := config.MustLoadConfig()

	log := logger.New(cfg.Env)

	collector := geonode.New(log)
	pager := collector.NewPageScheduler()

	for {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Collector.Timeout)
		defer cancel()

		pageUrl := pager()
		uuid := uuid.New().String()
		filename := fmt.Sprintf("storage/init/test/proxies_%s.json", uuid)

		log := log.With(slog.String("filename", filename), slog.String("page", pageUrl))

		_, err := collector.Collect(ctx, pageUrl, filename)
		if err != nil {
			log.Error("collector.Collect failed", slog.String("err", err.Error()))
			return
		}
		log.Info("collector done")
	}
}
