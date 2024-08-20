package favoritsservice

import (
	"errors"
	apiv1 "proxyfinder/internal/service/api"
	favoritsstorage "proxyfinder/internal/storage/sqlx/favorits"
	"proxyfinder/pkg/filter"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFavoritsService_ValidateOptions(t *testing.T) {
	type args struct {
		options filter.Options
	}
	type want struct {
		err error
	}
	type init struct {
		service *FavoritsService
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
			desc: "Valid data",
			init: func (t *testing.T) init {
				options := filter.New()
				options.SetLimit(10)
				options.SetOffset(0)
				options.AddField("user_id", filter.OpEq, "1", "int64")

				args := args{
					options: options,
				}
				want := want{
					err: nil,
				}

				return init{
					service: New(nil, favoritsstorage.New(nil)),
					args:    args,
					want:    want,
				}
			},
		},
		{
			name: "Positive test #2",
			desc: "Valid data",
			init: func (t *testing.T) init {
				options := filter.New()
				options.SetLimit(10)
				options.SetOffset(0)
				options.AddField("proxy_id", filter.OpEq, "1", "int64")

				args := args{
					options: options,
				}
				want := want{
					err: nil,
				}

				return init{
					service: New(nil, favoritsstorage.New(nil)),
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
			err := i.service.ValidateOptions(args.options)
			t.Logf("Execution time: %s", time.Since(start))
			require.Equal(t, want.err, err)
		})
	}
}

func TestFavoritsService_ValidateOptionsError(t *testing.T) {
	type args struct {
		options filter.Options
	}
	type want struct {
		err error
	}
	type init struct {
		service *FavoritsService
		args    args
		want    want
	}
	tests := []struct {
		name string
		desc string
		init func(t *testing.T) init
	}{
		{
			name: "Negative test #1",
			desc: "Invalid data",
			init: func (t *testing.T) init {
				options := filter.New()
				options.SetLimit(10)
				options.SetOffset(0)
				options.AddField("name", filter.OpEq, "serega", "string")

				args := args{
					options: options,
				}
				want := want{
					err: errors.New(apiv1.ErrInvalidField),
				}

				return init{
					service: New(nil, favoritsstorage.New(nil)),
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
			err := i.service.ValidateOptions(args.options)
			t.Logf("Execution time: %s", time.Since(start))
			require.Equal(t, want.err, err)
		})
	}
}
