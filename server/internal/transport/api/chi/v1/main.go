package chiapiv1

import (
	"context"
	"net/http"
	chiapi "proxyfinder/internal/transport/api/chi"
	"proxyfinder/pkg/options"
	"strconv"
	"strings"
)

const (
	FilterCtxKey     = "FilterOptions"
	SortCtxKey       = "SortOptions"
	ErrOptionsNotSet = "Options not set"
)

func FilterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			page, perPage int
		)
		opts := options.New()

		for param := range r.URL.Query() {
			// ignore sort params
			if param == "sort_by" || param == "sort_order" {
				continue
			}

			// if value is page or perPage, it's not a filter
			if param == "page" {
				page, _ = strconv.Atoi(r.URL.Query().Get(param))
				err := opts.AddField(param, options.OpEq, page)
				if err != nil {
					chiapi.JSONresponse(w, http.StatusBadRequest, nil, err)
					return
				}
				continue
			}
			if param == "perPage" {
				perPage, _ = strconv.Atoi(r.URL.Query().Get(param))
				err := opts.AddField(param, options.OpEq, perPage)
				if err != nil {
					chiapi.JSONresponse(w, http.StatusBadRequest, nil, err)
					return
				}
				continue
			}

			// if value has comma, it's not a single value
			if strings.Contains(r.URL.Query().Get(param), ",") {
				values := strings.Split(r.URL.Query().Get(param), ",")
				err := opts.AddField(param, options.OpIn, values)
				if err != nil {
					chiapi.JSONresponse(w, http.StatusBadRequest, nil, err)
					return
				}
				continue
			}	

			// Else it's a single value
			err := opts.AddField(param, options.OpEq, r.URL.Query().Get(param))
			if err != nil {
				chiapi.JSONresponse(w, http.StatusBadRequest, nil, err)
				return
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, FilterCtxKey, opts)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func SortMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sortList := strings.Split(r.URL.Query().Get("sort_by"), ",")
		orderList := strings.Split(r.URL.Query().Get("sort_order"), ",")

		opts := options.New()

		for i := range sortList {
			if i >= len(orderList) {
				opts.AddField(sortList[i], options.OpEq, "asc")
			} else {
				opts.AddField(sortList[i], options.OpEq, orderList[i])
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, SortCtxKey, opts)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
