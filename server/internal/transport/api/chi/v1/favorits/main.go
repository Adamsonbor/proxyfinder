package favoritsapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"proxyfinder/internal/domain/dto"
	apiv1 "proxyfinder/internal/service/api"
	serviceapiv1 "proxyfinder/internal/service/api"
	jwtservice "proxyfinder/internal/service/api/v1/auth/jwt"
	chiapi "proxyfinder/internal/transport/api/chi"
	chiapiv1 "proxyfinder/internal/transport/api/chi/v1"
	"proxyfinder/pkg/options"

	"github.com/go-chi/chi/v5"
)

type FavoritsController struct {
	log     *slog.Logger
	Router  *chi.Mux
	service apiv1.FavoritsService
	jwt     serviceapiv1.JWTService
}

func New(log *slog.Logger, service apiv1.FavoritsService, jwt serviceapiv1.JWTService) *FavoritsController {
	r := chi.NewRouter()
	router := &FavoritsController{
		log:     log,
		Router:  r,
		service: service,
		jwt:     jwt,
	}

	r.Use(chiapiv1.FilterMiddleware)
	r.Use(chiapiv1.SortMiddleware)
	r.Use(jwt.JWTMiddleware)
	r.Get("/", router.GetAll)
	r.Post("/", router.Create)
	r.Delete("/{proxy_id}", router.Delete)


	return router
}

// GetAll
// @Summary Get all favorits
// @Description Get all favorits with filters
// @Tags favorits
// @Produces json
// @Param Authorization header string true "Authorization"
// @Param page query int false "Page number"
// @Param perPage query int false "Results per page"
// @Param user_id query int false "User id"
// @Param proxy_id query int false "Proxy id"
// @Success 200 {object} chiapi.Response
// @Failure 400 {object} chiapi.Response
// @Failure 500 {object} chiapi.Response
// @Router /api/v1/favorits [get]
func (self *FavoritsController) GetAll(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "FavoritsController.GetAll"))

	filter, _ := r.Context().Value(chiapiv1.FilterCtxKey).(options.Options)
	sort, _ := r.Context().Value(chiapiv1.SortCtxKey).(options.Options)

	userId := r.Context().Value(jwtservice.JwtUserCtxKey).(int64)
	filter.AddField("user_id", options.OpEq, userId)

	log.Debug("request", slog.Any("filter", filter))
	log.Debug("request", slog.Any("sort", sort))

	result, err := self.service.GetAll(r.Context(), filter, sort)
	if err != nil {
		log.Error("failed to get all", slog.Any("err", err))
		chiapi.JSONresponse(w, http.StatusBadRequest, nil, err)
		return
	}

	chiapi.JSONresponse(w, http.StatusOK, result, nil)
}

// Create
// @Summary Create
// @Description Create favorit. Request must contain Authorization in headers and {user_id, proxy_id} in body
// @Tags favorits
// @Produces integer
// @Param Authorization header string true "Authorization"
// @Param body body dto.FavoriteCreate true "FavoriteCreate"
// @Success 200 {object} chiapi.Response
// @Failure 400 {object} chiapi.Response
// @Failure 500 {object} chiapi.Response
// @Router /api/v1/favorits/ [post]
func (self *FavoritsController) Create(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "FavoritsController.Create"))

	defer r.Body.Close()

	filter := options.New()

	var favorit dto.FavoriteCreate
	err := json.NewDecoder(r.Body).Decode(&favorit)
	if err != nil {
		log.Error("failed to decode", slog.Any("err", err))
		chiapi.JSONresponse(w, http.StatusBadRequest, nil, err)
		return
	}

	log.Debug("body", slog.Any("body", favorit))

	filter.AddField("user_id", options.OpEq, favorit.UserId)
	filter.AddField("proxy_id", options.OpEq, favorit.ProxyId)

	log.Debug("request", slog.Any("options", filter))

	id, err := self.service.Save(r.Context(), filter)
	if err != nil {
		log.Error("failed to create", slog.Any("err", err))
		chiapi.JSONresponse(w, http.StatusBadRequest, nil, err)
		return
	}

	chiapi.JSONresponse(w, http.StatusOK, id, nil)
}

// Delete
// @Summary Delete
// @Description Delete favorit. Request must contain proxy_id and access_token in headers for authorization
// @Tags favorits
// @Produces integer
// @Param proxy_id path string true "Proxy id"
// @Success 200 {object} chiapi.Response
// @Failure 400 {object} chiapi.Response
// @Failure 500 {object} chiapi.Response
// @Router /api/v1/favorits/{proxy_id} [delete]
func (self *FavoritsController) Delete(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "FavoritsController.Delete"))

	filter := options.New()

	proxyId := chi.URLParam(r, "proxy_id")
	userId := r.Context().Value(jwtservice.JwtUserCtxKey).(int64)

	filter.AddField("user_id", options.OpEq, userId)
	filter.AddField("proxy_id", options.OpEq, proxyId)

	log.Debug("request", slog.Any("options", filter))

	err := self.service.Delete(r.Context(), filter)
	if err != nil {
		log.Error("failed to delete", slog.Any("err", err))
		chiapi.JSONresponse(w, http.StatusBadRequest, nil, err)
		return
	}

	chiapi.JSONresponse(w, http.StatusOK, nil, nil)
}
