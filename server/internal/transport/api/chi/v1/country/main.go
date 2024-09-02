package countryapi

import (
	"log/slog"
	"net/http"
	serviceapiv1 "proxyfinder/internal/service/api"
	chiapi "proxyfinder/internal/transport/api/chi"
	chiapiv1 "proxyfinder/internal/transport/api/chi/v1"
	"proxyfinder/pkg/options"

	"github.com/go-chi/chi/v5"
)

type CountryController struct {
	log *slog.Logger
	Router *chi.Mux
	service serviceapiv1.CountryService
}

func New(log *slog.Logger, service serviceapiv1.CountryService) *CountryController {
	r := chi.NewRouter()

	router := &CountryController{
		log:     log,
		Router:  r,
		service: service,
	}

	r.Use(chiapiv1.FilterMiddleware)
	r.Use(chiapiv1.SortMiddleware)
	r.Get("/", router.GetAll)


	return router
}

// GetAll
// @Summary Get all countries
// @Description Get all countries with filters and sorting
// @Tags country
// @Produces json
// @Param page query int false "Page number"
// @Param perPage query int false "Results per page"
// @Param name query string false "Country name"
// @Param code query string false "Country code"
// @Success 200 {object} chiapi.Response
// @Failure 400 {object} chiapi.Response
// @Router /api/v1/country [get]
func (self *CountryController) GetAll(w http.ResponseWriter, r *http.Request) {
	log := self.log.With(slog.String("op", "CountryController.GetAll"))

	filter, _ := r.Context().Value(chiapiv1.FilterCtxKey).(options.Options)
	sort, _ := r.Context().Value(chiapiv1.SortCtxKey).(options.Options)

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
