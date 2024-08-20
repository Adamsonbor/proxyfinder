package chiapiv1

import (
	"context"
	"net/http"
	chiapi "proxyfinder/internal/transport/api/chi"
	"proxyfinder/pkg/filter"
	"proxyfinder/pkg/pagination"
	"strconv"
)

func FilterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		options := filter.New()

		for param := range r.URL.Query() {
			if param == "page" || param == "perPage" {
				continue
			}
			err := options.AddField(param, filter.OpEq, r.URL.Query().Get(param), "string")
			if err != nil {
				chiapi.JSONresponse(w, http.StatusBadRequest, nil, err)
				return
			}
		}

		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			page = 1
		}
		perPage, err := strconv.Atoi(r.URL.Query().Get("perPage"))
		if err != nil {
			perPage = 10
		}
		limit, offset := pagination.LimitOffset(page, perPage)
		options.SetLimit(limit)
		options.SetOffset(offset)

		ctx := r.Context()
		ctx = context.WithValue(ctx, filter.FilterCtxKey, options)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
