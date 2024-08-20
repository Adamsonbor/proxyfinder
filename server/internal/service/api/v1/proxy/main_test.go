package proxyservice

import (
	proxystorage "proxyfinder/internal/storage/sqlx/proxy"
	"proxyfinder/pkg/filter"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProxyService_OptionsMap(t *testing.T) {
	type want struct {
		options filter.Options
		err     error
	}
	type args struct {
		options filter.Options
	}
	type init struct {
		service ProxyService
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
				options := filter.New()
				country, status := "Russian Federation", "Unavailable"
				options.AddField("country_name", filter.OpEq, country, "string")
				options.AddField("status_name", filter.OpEq, status, "string")

				wantOpts := filter.New()
				wantOpts.AddField("country.name", filter.OpEq, country, "string")
				wantOpts.AddField("status.name", filter.OpEq, status, "string")

				return init{
					service: *New(nil, proxystorage.New(nil)),
					args:    args{options: options},
					want: want{
						options: wantOpts,
						err:     nil,
					},
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
			got, err := i.service.OptionsMap(args.options)
			t.Logf("Took %s", time.Since(start))
			require.NoError(t, err)
			assert.ElementsMatch(t, want.options.Fields(), got.Fields())
		})
	}
}

func TestProxyService_IsValudUpdateOptions(t *testing.T) {
	type args struct {
		options filter.Options
	}
	type want struct {
		err error
	}
	type init struct {
		service *ProxyService
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
				service := New(nil, proxystorage.New(nil))
				options := filter.New()
				options.AddField("id", filter.OpEq, 1, "int64")
				options.AddField("response_time", filter.OpEq, 1, "int64")
				options.AddField("status_id", filter.OpEq, 1, "int64")

				return init{
					service: service,
					args:    args{options: options},
					want: want{
						err: nil,
					},
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
			err := i.service.IsValudUpdateOptions(args.options)
			t.Logf("Execution took %s", time.Since(start))
			assert.Equal(t, want.err, err)
		})
	}
}
