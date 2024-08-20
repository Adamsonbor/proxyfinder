package favoritsapi

import (
	"log/slog"
	"net/http"
	apiv1 "proxyfinder/internal/service/api"
	chiapi "proxyfinder/internal/transport/api/chi"
	chiapiv1 "proxyfinder/internal/transport/api/chi/v1"
	"proxyfinder/pkg/filter"

	"github.com/go-chi/chi/v5"
)

type FavoritsController struct {
	log     *slog.Logger
	Router  *chi.Mux
	service apiv1.FavoritsService
}

func New(log *slog.Logger, service apiv1.FavoritsService) *FavoritsController {
	r := chi.NewRouter()
	router := &FavoritsController{
		log:     log,
		Router:  r,
		service: service,
	}

	r.Use(chiapiv1.FilterMiddleware)
	r.Get("/", router.GetAll)

	return router
}

// GetAll
// @Summary Get all favorits
// @Description Get all favorits with filters
// @Tags favorits
// @Produces json
// @Param page query int false "Page number"
// @Param perPage query int false "Results per page"
// @Param user_id query int false "User id"
// @Param proxy_id query int false "Proxy id"
// @Success 200 {object} chiapi.Response
// @Failure 400 {object} chiapi.Response
// @Failure 500 {object} chiapi.Response
// @Router /favorits [get]
func (self *FavoritsController) GetAll(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "FavoritsController.GetAll"))

	options, _ := r.Context().Value(filter.FilterCtxKey).(filter.Options)

	log.Debug("request", slog.Any("options", options))

	result, err := self.service.GetAll(r.Context(), options)
	if err != nil {
		log.Error("failed to get all", slog.Any("err", err))
		chiapi.JSONresponse(w, http.StatusBadRequest, nil, err)
		return
	}

	chiapi.JSONresponse(w, http.StatusOK, result, nil)
}
