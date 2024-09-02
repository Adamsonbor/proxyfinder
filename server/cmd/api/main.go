package main

import (
	"log/slog"
	"net/http"
	"proxyfinder/internal/config"
	googleservice "proxyfinder/internal/service/api/v1/auth/google"
	jwtservice "proxyfinder/internal/service/api/v1/auth/jwt"
	countryservice "proxyfinder/internal/service/api/v1/country"
	favoritsservice "proxyfinder/internal/service/api/v1/favorits"
	proxyservice "proxyfinder/internal/service/api/v1/proxy"
	userservice "proxyfinder/internal/service/api/v1/user"
	countrystorage "proxyfinder/internal/storage/sqlx/country"
	favoritsstorage "proxyfinder/internal/storage/sqlx/favorits"
	proxystorage "proxyfinder/internal/storage/sqlx/proxy"
	userstorage "proxyfinder/internal/storage/sqlx/user"
	googleapi "proxyfinder/internal/transport/api/chi/v1/auth/google"
	countryapi "proxyfinder/internal/transport/api/chi/v1/country"
	favoritsapi "proxyfinder/internal/transport/api/chi/v1/favorits"
	proxyapi "proxyfinder/internal/transport/api/chi/v1/proxy"
	userapi "proxyfinder/internal/transport/api/chi/v1/user"
	"proxyfinder/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	_ "proxyfinder/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Proxpro API
// @version 1.0
// @description Proxpro API
// @termsOfService http://swagger.io/terms/
// @contact.name Adamson Bor
// @contact.url http://github.com/Adamsonbor
// @contact.email adamsonbor@gmail.com
// @host localhost:8080
func main() {

	// INIT config
	cfg := config.MustLoadConfig()

	// INIT logger
	log := logger.New(cfg.Env)
	log.Info("Initializing with env: " + cfg.Env)

	// INIT storage
	db := sqlx.MustOpen("sqlite3", "./storage/local.db")
	proxyStorage := proxystorage.New(db)
	favoritsStorage := favoritsstorage.New(db)
	userStorage := userstorage.New(db)
	countryStorage := countrystorage.New(db)

	// INIT service
	proxyService := proxyservice.New(log, proxyStorage)
	favoritsService := favoritsservice.New(log, favoritsStorage)
	userService := userservice.New(log, userStorage)
	jwtService := jwtservice.New(log, cfg)
	googleAuthService := googleservice.New(log, userService, jwtService, cfg)
	countryService := countryservice.New(log, countryStorage)

	// INIT router
	mux := chi.NewRouter()
	proxyController := proxyapi.New(log, proxyService)
	favoritsController := favoritsapi.New(log, favoritsService, jwtService)
	googleAuthController := googleapi.New(log, googleAuthService, *cfg)
	userController := userapi.New(log, userService, jwtService, *cfg)
	countryController := countryapi.New(log, countryService)

	// register routes
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Heartbeat("/ping"))
	mux.Use(middleware.Heartbeat("/health"))
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Route("/api/v1", func(r chi.Router) {
		r.Mount("/proxy", proxyController.Router)
		r.Mount("/favorits", favoritsController.Router)
		r.Mount("/user", userController.Router)
		r.Mount("/country", countryController.Router)
	})
	mux.Mount("/auth/google", googleAuthController.Router)
	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// print routes
	chi.Walk(mux, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Info(route, slog.String("method", method))
		return nil
	})

	// run server
	http.ListenAndServe(":8080", mux)
}
