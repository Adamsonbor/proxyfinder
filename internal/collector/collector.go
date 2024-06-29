package collector

import (
	"context"
	"proxyfinder/internal/domain"
)

// Collect the data to store in db
type Collector interface {
	Collect(ctx context.Context) ([]domain.Proxy, error)
}
