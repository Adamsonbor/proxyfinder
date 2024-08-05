package rabbit

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"proxyfinder/internal/api"
	"proxyfinder/internal/broker/rabbit"
	"proxyfinder/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Router struct {
	log    *slog.Logger
	Router *chi.Mux
	rabbit *rabbit.RabbitService
	cfg    *config.Config
}

func New(log *slog.Logger, rabbit *rabbit.RabbitService, cfg *config.Config) *Router {
	r := chi.NewRouter()

	router := &Router{
		log:    log,
		Router: r,
		rabbit: rabbit,
		cfg:    cfg,
	}

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		ExposedHeaders: []string{"Content-Range"},
		MaxAge:         300, // Maximum value not ignored by any of major browsers
	}))


	r.Route("/", func(r chi.Router) {
		r.Post("/publish", router.Publish)
	})

	return router
}

func (ro *Router) Publish(w http.ResponseWriter, r *http.Request) {
	log := ro.log.With(slog.String("op", "rabbit.Publish"))
	log.Debug("Publish request")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Warn("io.ReadAll failed", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.Response{Status: "error", Error: err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), ro.cfg.Rabbit.Timeout)
	defer cancel()

	err = ro.rabbit.Publish(ctx, body)
	if err != nil {
		log.Warn("rabbit.Publish failed", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.Response{Status: "error", Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(api.Response{Status: "success"})

	log.Debug("Publish success")
}
