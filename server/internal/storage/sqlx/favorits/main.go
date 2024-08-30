package favoritsstorage

import (
	"context"
	"proxyfinder/internal/domain"
	apiv1 "proxyfinder/internal/service/api"
	"proxyfinder/internal/storage"
	"proxyfinder/pkg/options"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	selectAllQuery = `SELECT * FROM favorits`
	insertQuery    = `INSERT INTO favorits (user_id, proxy_id, created_at, updated_at) VALUES (?, ?, ?, ?)`
	deleteQuery    = `DELETE FROM favorits WHERE user_id = ? AND proxy_id = ?`
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

func (self *FavoritsStorage) Save(ctx context.Context, opts options.Options) (int64, error) {
	opts.AddField("created_at", options.OpEq, time.Now().Unix())
	opts.AddField("updated_at", options.OpEq, time.Now().Unix())

	res, err := self.db.ExecContext(ctx, insertQuery, opts.Values()...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (self *FavoritsStorage) Delete(ctx context.Context, opts options.Options) error {
	_, err := self.db.ExecContext(ctx, deleteQuery, opts.Values()...)
	return err
}
