package sqlite

import (
	"context"
	"database/sql"
	"proxyfinder/internal/domain"
	"proxyfinder/internal/storage"

	"github.com/jmoiron/sqlx"
)

type ProxyStorage struct {
	db *sqlx.DB
}

func NewProxy(db *sqlx.DB) *ProxyStorage {
	return &ProxyStorage{
		db: db,
	}
}

func (s *ProxyStorage) BeginTx(ctx context.Context, opt *sql.TxOptions) (*sqlx.Tx, error) {
	return s.db.BeginTxx(ctx, opt)
}

func (s *ProxyStorage) Begin() (*sqlx.Tx, error) {
	return s.db.Beginx()
}

func (s *ProxyStorage) Update(ctx context.Context, tx *sqlx.Tx, id int64, fields *storage.ProxyUpdate) error {
	var proxy domain.Proxy
	err := tx.GetContext(ctx, &proxy, "SELECT * FROM proxy WHERE id = ?", id)
	if err != nil {
		return err
	}

	if fields.Ip != nil {
		proxy.Ip = *fields.Ip
	}
	if fields.Port != nil {
		proxy.Port = *fields.Port
	}
	if fields.Protocol != nil {
		proxy.Protocol = *fields.Protocol
	}
	if fields.ResponseTime != nil {
		proxy.ResponseTime = *fields.ResponseTime
	}
	if fields.CountryId != nil {
		proxy.CountryId = *fields.CountryId
	}
	if fields.StatusId != nil {
		proxy.StatusId = *fields.StatusId
	}

	query := `UPDATE proxy
		SET ip = ?, port = ?,
		protocol = ?, country_id = ?,
		status_id = ?, response_time = ?,
		updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err = tx.ExecContext(
		ctx, query, proxy.Ip, proxy.Port,
		proxy.Protocol, proxy.CountryId,
		proxy.StatusId, proxy.ResponseTime, id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProxyStorage) Save(ctx context.Context, proxy *domain.Proxy) (int64, error) {
	query := "INSERT INTO proxy (ip, port, protocol, status_id, country_id) VALUES (?, ?, ?, ?, ?) RETURNING id"

	res, err := s.db.ExecContext(ctx, query, proxy.Ip, proxy.Port, proxy.Protocol, proxy.StatusId, proxy.CountryId)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (s *ProxyStorage) SaveAll(ctx context.Context, proxies []domain.Proxy) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO proxy (ip, port, protocol, status_id, country_id)
							VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range proxies {
		_, err := stmt.Exec(v.Ip, v.Port, v.Protocol, v.StatusId, v.CountryId)
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

func (s *ProxyStorage) GetAvailable(ctx context.Context) ([]domain.Proxy, error) {
	query := "SELECT * FROM proxy WHERE status_id = 2"

	var proxies []domain.Proxy
	err := s.db.SelectContext(ctx, &proxies, query)

	return proxies, err
}

func (s *ProxyStorage) GetAll(ctx context.Context) ([]domain.Proxy, error) {
	query := "SELECT * FROM proxy"

	var proxies []domain.Proxy
	err := s.db.SelectContext(ctx, &proxies, query)

	return proxies, err
}

func (s *ProxyStorage) UpdateStatus(ctx context.Context, id int64, statusId int64) error {
	query := "UPDATE proxy SET status_id = ? WHERE id = ?"

	_, err := s.db.ExecContext(ctx, query, statusId, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProxyStorage) Deletex(ctx context.Context, tx *sqlx.Tx, id int64) error {
	query := "DELETE FROM proxy WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, id)
	return err
}

func (s *ProxyStorage) Savex(ctx context.Context, tx *sqlx.Tx, inst *domain.Proxy) (int64, error) {
	query := "INSERT INTO proxy (ip, port, protocol, status_id, country_id) VALUES (?, ?, ?, ?, ?) RETURNING id"
	res, err := tx.ExecContext(ctx, query, inst.Ip, inst.Port, inst.Protocol, inst.StatusId, inst.CountryId)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (s *ProxyStorage) SaveAllx(ctx context.Context, tx *sqlx.Tx, proxies []domain.Proxy) error {
	stmt, err := tx.Prepare(`INSERT INTO proxy (ip, port, protocol, status_id, country_id)
							VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range proxies {
		_, err := stmt.Exec(v.Ip, v.Port, v.Protocol, v.StatusId, v.CountryId)
		if err != nil {
			return err
		}
	}

	return nil
}
