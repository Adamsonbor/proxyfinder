package countrystorage

import (
	"context"
	"proxyfinder/internal/domain"
	"proxyfinder/internal/storage"
	"proxyfinder/pkg/options"

	"github.com/jmoiron/sqlx"
)

const (
	getAllQuery = `SELECT * FROM country`
)

type CountryStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *CountryStorage {
	return &CountryStorage{
		db: db,
	}
}

func (self *CountryStorage) GetAll(
	ctx context.Context,
	filter options.Options,
	sort options.Options,
) ([]domain.Country, error) {
	qb := storage.NewQueryBuilder()
	err := qb.Filter(filter)
	if err != nil {
		return nil, err
	}

	qb.Sort(sort)

	var res []domain.Country
	err = self.db.SelectContext(ctx, &res, qb.BuildQuery(getAllQuery), qb.Values()...)
	return res, err
}
