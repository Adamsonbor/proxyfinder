package proxystorage

import (
	"context"
	"proxyfinder/internal/domain/dto"
	"proxyfinder/pkg/filter"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProxyStorage_GetAll(t *testing.T) {
	type want struct {
		proxies []dto.Proxy
		err     error
	}
	type args struct {
		ctx     context.Context
		options filter.Options
	}
	type init struct {
		storage ProxyStorage
		args    args
		want    want
	}
	tests := []struct {
		name string
		desc string
		init func(t *testing.T) init
	}{
		{
			name: "Positive test #1",
			desc: "Test with correct data",
			init: func(t *testing.T) init {
				db := sqlx.MustConnect("sqlite3", "../../../../storage/test.db")
				storage := New(db)

				query := selectWithStatusAndCountryQuery + `
					WHERE country.name LIKE ? AND status.name LIKE ? LIMIT ? OFFSET ?`

				country, status := "Russian Federation", "Unavailable"
				limit, offset := 10, 0
				var answerProxies []dto.Proxy
				err := db.Select(&answerProxies, query, country, status, limit, offset)
				require.NoError(t, err)

				options := filter.New()
				options.SetLimit(limit)
				options.SetOffset(offset)
				options.AddField("country.name", filter.OpLike, country, "string")
				options.AddField("status.name", filter.OpLike, status, "string")

				args := args{
					ctx:     context.Background(),
					options: options,
				}
				want := want{
					proxies: answerProxies,
					err:     nil,
				}

				return init{
					storage: storage,
					args:    args,
					want:    want,
				}
			},
		},
		{
			name: "Positive test #2",
			desc: "Test with correct data and paginations",
			init: func(t *testing.T) init {
				db := sqlx.MustConnect("sqlite3", "../../../../storage/test.db")
				storage := New(db)

				query := selectWithStatusAndCountryQuery + `
					WHERE country.name LIKE ? AND status.name LIKE ? LIMIT ? OFFSET ?`

				country, status := "Russian Federation", "Unavailable"
				limit, offset := 20, 10
				var answerProxies []dto.Proxy
				err := db.Select(&answerProxies, query, country, status, limit, offset)
				require.NoError(t, err)

				options := filter.New()
				options.SetLimit(limit)
				options.SetOffset(offset)
				options.AddField("country.name", filter.OpLike, country, "string")
				options.AddField("status.name", filter.OpLike, status, "string")

				args := args{
					ctx:     context.Background(),
					options: options,
				}
				want := want{
					proxies: answerProxies,
					err:     nil,
				}

				return init{
					storage: storage,
					args:    args,
					want:    want,
				}
			},
		},
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := tt.init(t)
			args := i.args
			want := i.want

			start := time.Now()
			proxies, err := i.storage.GetAll(args.ctx, args.options)
			t.Logf("GetAll took %s", time.Since(start))
			assert.Equal(t, want.err, err)
			assert.ElementsMatch(t, want.proxies, proxies)
			// t.Logf("expected len: %d got: %d", len(want.proxies), len(proxies))
		})
	}
}
