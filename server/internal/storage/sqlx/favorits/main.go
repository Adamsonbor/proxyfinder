package favoritsstorage

import (
	"context"
	"proxyfinder/internal/domain"
	apiv1 "proxyfinder/internal/service/api"
	"proxyfinder/internal/storage"
	"proxyfinder/pkg/options"

	"github.com/jmoiron/sqlx"
)

const (
	selectAllQuery = `SELECT * FROM favorits`
	ErrInvalidType = "Invalid type"
)

type FavoritsStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) apiv1.FavoritsStorage {
	return &FavoritsStorage{
		db: db,
	}
}

func (self *FavoritsStorage) GetAll(
	ctx context.Context,
	filter options.Options,
	sort options.Options,
) ([]domain.Favorits, error) {
	qb := storage.NewQueryBuilder()
	err := qb.Filter(filter)
	if err != nil {
		return nil, err
	}
	qb.Sort(sort)

	var res []domain.Favorits
	err = self.db.SelectContext(ctx, &res, qb.BuildQuery(selectAllQuery), qb.Values()...)
	return res, err
}
