package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"proxyfinder/internal/domain"
	"proxyfinder/internal/storage/dto"

	"github.com/jmoiron/sqlx"
)

var (
	ErrRecordNotFound = fmt.Errorf("Recod not found")
	ErrEmptyOptions   = fmt.Errorf("Empty options")
	ErrInvalidId      = fmt.Errorf("Invalid id")
)

type ProxyRepo interface {
	GetAll(ctx context.Context, page, perPage int, country, status string) ([]dto.ProxyDTO, error)
}

type UserStorage interface {
	Begin(ctx context.Context, opt *sql.TxOptions) (*sqlx.Tx, error)

	Get(ctx context.Context, id int64) (*domain.User, error)
	GetBy(ctx context.Context, o map[string]interface{}) ([]domain.User, error)
	Create(ctx context.Context, tx *sqlx.Tx, inst *domain.User) (int64, error)
	Update(ctx context.Context, tx *sqlx.Tx, id int64, o map[string]interface{}) error
	Delete(ctx context.Context, tx *sqlx.Tx, id int64) error
}

type SessionStorage interface {
	Begin(ctx context.Context, opt *sql.TxOptions) (*sqlx.Tx, error)

	Get(ctx context.Context, id int64) (*domain.Session, error)
	GetBy(ctx context.Context, o map[string]interface{}) ([]domain.Session, error)
	Create(ctx context.Context, tx *sqlx.Tx, inst *domain.Session) (int64, error)
	Update(ctx context.Context, tx *sqlx.Tx, id int64, o map[string]interface{}) error
	Delete(ctx context.Context, tx *sqlx.Tx, id int64) error
}

type ProxyStorage interface {
	Begin(ctx context.Context, opt *sql.TxOptions) (*sqlx.Tx, error)

	Get(ctx context.Context, id int64) (*dto.ProxyDTO, error)
	GetBy(ctx context.Context, o map[string]interface{}) ([]dto.ProxyDTO, error)
	Create(ctx context.Context, tx *sqlx.Tx, inst *domain.Proxy) (int64, error)
	Update(ctx context.Context, tx *sqlx.Tx, id int64, o map[string]interface{}) error
	Delete(ctx context.Context, tx *sqlx.Tx, id int64) error
}

type CountryStorage interface {
	Begin(ctx context.Context, opt *sql.TxOptions) (*sqlx.Tx, error)

	Get(ctx context.Context, id int64) (*domain.Country, error)
	GetBy(ctx context.Context, o map[string]interface{}) ([]domain.Country, error)
	Create(ctx context.Context, tx *sqlx.Tx, inst *domain.Country) (int64, error)
	Update(ctx context.Context, tx *sqlx.Tx, id int64, o map[string]interface{}) error
	Delete(ctx context.Context, tx *sqlx.Tx, id int64) error
}

type StatusStorage interface {
	Begin(ctx context.Context, opt *sql.TxOptions) (*sqlx.Tx, error)

	Get(ctx context.Context, id int64) (*domain.Status, error)
	GetBy(ctx context.Context, o map[string]interface{}) ([]domain.Status, error)
	Create(ctx context.Context, tx *sqlx.Tx, inst *domain.Status) (int64, error)
	Update(ctx context.Context, tx *sqlx.Tx, id int64, o map[string]interface{}) error
	Delete(ctx context.Context, tx *sqlx.Tx, id int64) error
}

type FavoritsStorage interface {
	Begin(ctx context.Context, opt *sql.TxOptions) (*sqlx.Tx, error)

	Get(ctx context.Context, id int64) (*domain.Favorits, error)
	GetBy(ctx context.Context, o map[string]interface{}) ([]domain.Favorits, error)
	Create(ctx context.Context, tx *sqlx.Tx, inst *domain.Favorits) (int64, error)
	Update(ctx context.Context, tx *sqlx.Tx, id int64, o map[string]interface{}) error
	Delete(ctx context.Context, tx *sqlx.Tx, id int64) error
}

func ErrRecordNotFoundWrap(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrRecordNotFound
	}
	return err
}
