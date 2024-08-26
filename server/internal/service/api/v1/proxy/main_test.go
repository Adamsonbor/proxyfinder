package proxyservice

import (
	proxystorage "proxyfinder/internal/storage/sqlx/proxy"
	opts "proxyfinder/pkg/options"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProxyService_OptionsMap(t *testing.T) {
	type want struct {
		options opts.Options
		err     error
	}
	type args struct {
		options opts.Options
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
				options := opts.New()
				country, status := "Russian Federation", "Unavailable"
				options.AddField("country_name", opts.OpEq, country)
				options.AddField("status_name", opts.OpEq, status)

				wantOpts := opts.New()
				wantOpts.AddField("country.name", opts.OpEq, country)
				wantOpts.AddField("status.name", opts.OpEq, status)

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
			err := i.service.FieldsMap(args.options)
			t.Logf("Took %s", time.Since(start))
			require.NoError(t, err)
			assert.ElementsMatch(t, want.options.Fields(), args.options.Fields())
		})
	}
}

func TestProxyService_IsValudUpdateOptions(t *testing.T) {
	type args struct {
		options opts.Options
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
				options := opts.New()
				options.AddField("id", opts.OpEq, 1)
				options.AddField("response_time", opts.OpEq, 1)
				options.AddField("status_id", opts.OpEq, 1)

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
