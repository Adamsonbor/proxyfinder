package favoritsstorage

import (
	"bytes"
	"context"
	"fmt"
	"proxyfinder/internal/domain"
	apiv1 "proxyfinder/internal/service/api"
	"proxyfinder/pkg/filter"

	"github.com/jmoiron/sqlx"
)

const (
	selectAllQuery = `SELECT * FROM favorits`
)

type FavoritsStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) apiv1.FavoritsStorage {
	return &FavoritsStorage{
		db: db,
	}
}

func (self *FavoritsStorage) GetAll(ctx context.Context, options filter.Options) ([]domain.Favorits, error) {
	buf := bytes.Buffer{}
	buf.WriteString(selectAllQuery)

	if options.Is() {
		buf.WriteString(" WHERE ")
	}
	for i, option := range options.Fields() {
		if i != 0 {
			buf.WriteString(" AND ")
		}
		buf.WriteString(fmt.Sprintf("%s %s ?", option.Name, option.Op))
	}
	limit, offset := options.Limit(), options.Offset()
	values := options.Values()
	if limit > 0 {
		buf.WriteString(" LIMIT ?")
		values = append(values, limit)
	}
	if offset > 0 {
		buf.WriteString(" OFFSET ?")
		values = append(values, offset)
	}

	var res []domain.Favorits
	err := self.db.SelectContext(ctx, &res, buf.String(), values...)

	return res, err
}
