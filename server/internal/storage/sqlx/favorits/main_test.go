package favoritsstorage

import (
	"context"
	"proxyfinder/internal/domain"
	apiv1 "proxyfinder/internal/service/api"
	"proxyfinder/pkg/options"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFavoritsStorage_GetAll(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts options.Options
	}
	type want struct {
		res []domain.Favorits
		err error
	}
	type init struct {
		storage apiv1.FavoritsStorage
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
			desc: "Test with correct input",
			init: func(t *testing.T) init {
				db := sqlx.MustOpen("sqlite3", "../../../../storage/test.db")
				query := "INSERT INTO favorits (user_id, proxy_id) VALUES (?, ?)"
				user_id := 1
				_, err := db.Exec(query, user_id, 12)
				if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
					err = nil
				}
				require.NoError(t, err)
				_, err = db.Exec(query, user_id, 13)
				if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
					err = nil
				}
				answer := []domain.Favorits{}
				err = db.Select(&answer, "SELECT * FROM favorits WHERE user_id = ?", user_id)

				storage := New(db)

				opts := options.New()
				opts.AddField("user_id", "=", strconv.Itoa(user_id))
				args := args{
					ctx:     context.Background(),
					opts: opts,
				}
				want := want{
					res: answer,
					err: nil,
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
			got, err := i.storage.GetAll(args.ctx, args.opts, options.New())
			t.Logf("Execution time: %s", time.Since(start))
			require.Equal(t, want.err, err)
			assert.ElementsMatch(t, want.res, got)
		})
	}
}
