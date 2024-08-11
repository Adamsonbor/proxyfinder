package router

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"proxyfinder/internal/api"
	"proxyfinder/internal/auth"
	"proxyfinder/internal/domain"
	"proxyfinder/internal/storage"
	gormstorage "proxyfinder/internal/storage/v2/gorm-sotrage"
	sqlxstorage "proxyfinder/internal/storage/v2/sqlx-storage"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	log          *slog.Logger
	Router       *chi.Mux
	storage      *gormstorage.Storage
	proxyStorage *sqlxstorage.ProxyStorage
	jwt          auth.JWTService
}

func New(
	log *slog.Logger,
	storage *gormstorage.Storage,
	proxyStorage *sqlxstorage.ProxyStorage,
	jwt auth.JWTService,
) *Server {
	r := chi.NewRouter()
	s := Server{
		log:          log,
		Router:       r,
		storage:      storage,
		proxyStorage: proxyStorage,
		jwt:          jwt,
	}

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Content-Range"},
	}))

	NewUserRouter(log, r, storage, jwt)
	NewFavoritsRouter(log, r, storage, jwt)

	r.Route("/status", s.Crud)
	r.Route("/country", s.Crud)
	r.Route("/proxy", func(r chi.Router) {
		r.Route("/{id}", func(r chi.Router) {
			r.Use(idPermissionMiddleware)
			r.Get("/", s.Get)
			r.Put("/", s.Update)
			r.Delete("/", s.Delete)
		})
		r.Get("/", s.GetAllProxy)
		r.Post("/", s.Create)
	})

	return &s
}

func (s *Server) Crud(r chi.Router) {
	r.Get("/", s.GetAll)
	r.Post("/", s.Create)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(idPermissionMiddleware)
		r.Get("/", s.Get)
		r.Put("/", s.Update)
		r.Delete("/", s.Delete)
	})
}

func (s *Server) GetAll(w http.ResponseWriter, r *http.Request) {
	insts := s.GetSliceOfType(r.URL.Path)
	if insts == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := s.storage.GetAllFilter(insts, s.GetParams(r))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := api.Response{Status: "success", Data: insts}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	log := s.log.With("path", r.URL.Path, "method", r.Method)

	inst := s.GetType(r.URL.Path)
	if inst == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userId := r.Context().Value("id").(int64)

	log.Info("params", slog.Int64("id", userId))

	err := s.storage.Get(inst, userId)
	if err != nil {
		log.Error("get", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := api.Response{Status: "success", Data: inst}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	log := s.log.With("path", r.URL.Path, "method", r.Method)

	inst := s.GetType(r.URL.Path)
	if inst == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&inst)
	if err != nil {
		log.Error("decode", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inst, err = s.storage.Create(inst)
	if err != nil {
		log.Error("create", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(inst)
	if err != nil {
		log.Error("encode", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) Update(w http.ResponseWriter, r *http.Request) {
	log := s.log.With("path", r.URL.Path, "method", r.Method)

	inst := s.GetType(r.URL.Path)
	if inst == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&inst)
	if err != nil {
		log.Error("decode", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId := r.Context().Value("id").(int64)

	log.Info("params", slog.Int64("id", userId))

	defer func() {
		if recover() != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	reflect.ValueOf(inst).Elem().Field(0).SetInt(userId)
	insts, err := s.storage.Update(inst)
	if err != nil {
		log.Error("update", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(insts)
	if err != nil {
		log.Error("encode", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	log := s.log.With("path", r.URL.Path, "method", r.Method)

	userId := r.Context().Value("id").(int64)

	log.Info("params", slog.Int64("id", userId))

	inst := s.GetType(r.URL.Path)
	if inst == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer func() {
		if recover() != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	reflect.ValueOf(inst).Elem().Field(0).SetInt(userId)
	err := s.storage.Delete(inst)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetAllProxy(w http.ResponseWriter, r *http.Request) {
	log := s.log.With("path", r.URL.Path, "method", r.Method)
	log.Info("Start")

	page, perPage := 1, 10

	query := r.URL.Query()
	var err error
	if query.Get("page") != "" {
		page, err = strconv.Atoi(query.Get("page"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if query.Get("perPage") != "" {
		perPage, err = strconv.Atoi(query.Get("perPage"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	proxies, err := s.proxyStorage.GetAll(context.Background(), &storage.Options{Page: page, PerPage: perPage})
	if err != nil {
		log.Error("get all proxies", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := api.Response{Status: "success", Data: proxies}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Error("encode", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Info("Done")
}

func (self *Server) GetParams(r *http.Request) *storage.Options {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	perPage, _ := strconv.Atoi(query.Get("perPage"))
	return &storage.Options{
		Page:    page,
		PerPage: perPage,
	}
}

// Check if id is a number and set it in context
func idPermissionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		userId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "id", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) GetType(path string) interface{} {
	if strings.Contains(path, "proxy") {
		return &domain.Proxy{}
	} else if strings.Contains(path, "status") {
		return &domain.Status{}
	} else if strings.Contains(path, "country") {
		return &domain.Country{}
	}
	return nil
}

func (s *Server) GetSliceOfType(path string) interface{} {
	if strings.Contains(path, "proxy") {
		return &[]domain.Proxy{}
	} else if strings.Contains(path, "status") {
		return &[]domain.Status{}
	} else if strings.Contains(path, "country") {
		return &[]domain.Country{}
	}
	return nil
}

func JSONResponse(
	w http.ResponseWriter,
	status string,
	data interface{},
	err error,
) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	res := api.Response{Status: status, Data: data, Error: errMsg}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}

func ReturnJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
