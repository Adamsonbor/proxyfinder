package checker

import (
	"context"
	"fmt"
	"proxyfinder/internal/domain"
)

var (
	ErrInvalidProxy = fmt.Errorf("Invalid proxy configuration")
)

type Checker interface {
	Check(ctx context.Context, inst *domain.Proxy) (bool, error)
}
