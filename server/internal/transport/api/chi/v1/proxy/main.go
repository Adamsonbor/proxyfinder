package proxyapi

import (
	"fmt"
	"log/slog"
	"net/http"
	serviceapi "proxyfinder/internal/service/api"
	chiapi "proxyfinder/internal/transport/api/chi"
	apiv1 "proxyfinder/internal/transport/api/chi/v1"
	"proxyfinder/pkg/filter"

	"github.com/go-chi/chi/v5"
)

type ProxyController struct {
	log     *slog.Logger
	Router  *chi.Mux
	service serviceapi.ProxyService
}

func New(log *slog.Logger, service serviceapi.ProxyService) *ProxyController {
	r := chi.NewRouter()
	router := &ProxyController{
		log:     log,
		Router:  r,
		service: service,
	}

	r.Use(apiv1.FilterMiddleware)
	r.Get("/", router.GetAll)

	return router
}

// GetAll
// @Summary Get all proxies
// @Description Get all proxies
// @Tags proxy
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param perPage query int false "Number of items per page"
// @Param country_name query string false "Country name"
// @Param country_code query string false "Country code"
// @Param status_name query string false "Status name"
// @Param ip query string false "IP"
// @Param port query int false "Port"
// @Param protocol query string false "Protocol"
// @Param response_time query int false "Response time"
// @Success 200 {object} chiapi.Response
// @Error 400 {object} chiapi.Response
// @Error 500 {object} chiapi.Response
// @Router /proxy [get]
func (self *ProxyController) GetAll(w http.ResponseWriter, r *http.Request) {
	options, ok := r.Context().Value(filter.FilterCtxKey).(filter.Options)
	if !ok {
		chiapi.JSONresponse(w, http.StatusInternalServerError, nil, fmt.Errorf("failed to get options from context"))
		return
	}

	proxies, err := self.service.GetAll(r.Context(), options)
	if err != nil {
		chiapi.JSONresponse(w, http.StatusInternalServerError, nil, err)
		return
	}

	chiapi.JSONresponse(w, http.StatusOK, proxies, nil)
}
