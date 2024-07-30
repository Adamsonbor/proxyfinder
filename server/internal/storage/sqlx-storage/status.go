package sqlxstorage

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Status struct {
	db *sqlx.DB
}

func NewStatus(db *sqlx.DB) *Status {
	return &Status{db: db}
}

func (c *Status) SaveAll(ctx context.Context, names []string) error {
	tx, err := c.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	query := "INSERT INTO status (name) VALUES (?)"
	stmt, err := tx.Preparex(query)
	if err != nil {
		return err
	}

	for i := range names {
		_, err := stmt.ExecContext(ctx, query, names[i])
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
