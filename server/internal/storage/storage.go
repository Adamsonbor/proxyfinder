package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"proxyfinder/internal/domain"

	"github.com/jmoiron/sqlx"
)

var (
	ErrRecordNotFound = fmt.Errorf("Recod not found")
)

type Options struct {
	PerPage int
	Page    int
}

type ProxyUpdate struct {
	Ip           *string
	Port         *int
	Protocol     *string
	ResponseTime *int64
	StatusId     *int64
	CountryId    *int64
}

type UserStorage interface {
	Begin(ctx context.Context) (*sqlx.Tx, error)
	GetBy(ctx context.Context, field string, value interface{}) (*domain.User, error)
	Create(ctx context.Context, tx *sqlx.Tx, user *domain.User) (*domain.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.User, error)
	UpdateSession(ctx context.Context, tx *sqlx.Tx, user_id int64, refreshToken string) error
	NewSession(ctx context.Context, tx *sqlx.Tx, user_id int64, refreshToken string) error
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

func ErrRecordNotFoundWrap(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrRecordNotFound
	}
	return err
}
