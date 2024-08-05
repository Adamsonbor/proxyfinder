package googleapi

import (
	"log/slog"
	googleauth "proxyfinder/internal/auth/google"
	"proxyfinder/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Router struct {
	log    *slog.Logger
	cfg    *config.Config
	Router *chi.Mux
	gAuth  *googleauth.GoogleAuth
}

func NewRouter(log *slog.Logger, cfg *config.Config, gAuth *googleauth.GoogleAuth) *Router {
	r := chi.NewRouter()

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
		r.Use(gAuth.TokenOrLogin)
		r.Get("/login", gAuth.Login)
		r.Get("/logout", gAuth.Logout)
	})
	r.Get("/callback", gAuth.Callback)

	router := &Router{
		log:    log,
		cfg:    cfg,
		Router: r,
		gAuth:  gAuth,
	}

	return router
}
