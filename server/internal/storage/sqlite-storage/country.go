package sqlite

import (
	"context"
	"proxyfinder/internal/domain"

	"github.com/jmoiron/sqlx"
)

type Country struct {
	db *sqlx.DB
}

func NewCountry(db *sqlx.DB) *Country {
	return &Country{db: db}
}

func (c *Country) Save(ctx context.Context, inst *domain.Country) (int64, error) {
	query := "INSERT INTO country (name, code) VALUES (?, ?) RETURNING id"

	res, err := c.db.ExecContext(ctx, query, inst.Name, inst.Code)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (c *Country) GetByCode(ctx context.Context, code string) (*domain.Country, error) {
	var country domain.Country

	query := "SELECT * FROM country WHERE code = ?"

	err := c.db.SelectContext(ctx, &country, query)
	if err != nil {
		return nil, err
	}

	return &country, nil
}

func (c *Country) GetAll(ctx context.Context) ([]domain.Country, error) {
	var countries []domain.Country

	query := "SELECT * FROM country"
	err := c.db.SelectContext(ctx, &countries, query)
	if err != nil {
		return nil, err
	}

	return countries, nil
}

func (c *Country) SaveAll(ctx context.Context, insts []domain.Country) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := "INSERT INTO country (name, code) VALUES (?, ?)"
	stmt, err := tx.Preparex(query)
	if err != nil {
		return err
	}

	for i := range insts {
		_, err := stmt.ExecContext(ctx, query, insts[i].Name, insts[i].Code)
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

func (c *Country) Savex(ctx context.Context, tx *sqlx.Tx, inst *domain.Country) (int64, error) {
	query := "INSERT INTO country (name, code) VALUES (?, ?) RETURNING id"

	res, err := tx.ExecContext(ctx, query, inst.Name, inst.Code)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (c *Country) SaveAllx(ctx context.Context, tx *sqlx.Tx, insts []domain.Country) error {
	query := "INSERT INTO country (name, code) VALUES (?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := range insts {
		_, err := stmt.ExecContext(ctx, query, insts[i].Name, insts[i].Code)
		if err != nil {
			return err
		}
	}

	return nil
}
