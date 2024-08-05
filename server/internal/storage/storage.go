package storage

import (
	"context"
	"proxyfinder/internal/domain"

	"github.com/jmoiron/sqlx"
)

type ProxyUpdate struct {
	Ip           *string
	Port         *int
	Protocol     *string
	ResponseTime *int64
	StatusId     *int64
	CountryId    *int64
}

type ProxyStorage interface {
	Get(ctx context.Context, id int64) (*domain.Proxy, error)
	GetAvailable(ctx context.Context) ([]domain.Proxy, error)
	GetAll(ctx context.Context) ([]domain.Proxy, error)
	Update(ctx context.Context, tx *sqlx.Tx, id int64, fields *ProxyUpdate) error
	UpdateStatus(ctx context.Context, id int64, statusId int64) error
	Save(ctx context.Context, proxy *domain.Proxy) (int64, error)
	SaveAll(ctx context.Context, proxies []domain.Proxy) error
	Savex(ctx context.Context, tx *sqlx.Tx, inst *domain.Proxy) (int64, error)
	SaveAllx(ctx context.Context, tx *sqlx.Tx, insts []domain.Proxy) error

	Begin() (*sqlx.Tx, error)
}

type CountryStorage interface {
	Save(ctx context.Context, inst *domain.Country) (int64, error)
	GetByCode(ctx context.Context, code string) (*domain.Country, error)
	GetAll(ctx context.Context) ([]domain.Country, error)
	SaveAll(ctx context.Context, insts []domain.Country) error
	Savex(ctx context.Context, tx *sqlx.Tx, inst *domain.Country) (int64, error)
	SaveAllx(ctx context.Context, tx *sqlx.Tx, insts []domain.Country) error
}

type StatusStorage interface {
	SaveAll(ctx context.Context, insts []domain.Status) error
}
