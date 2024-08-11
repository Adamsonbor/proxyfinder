package sqlxstorage

import (
	"context"
	"proxyfinder/internal/storage"
	"proxyfinder/internal/storage/v2/dto"

	"github.com/jmoiron/sqlx"
)

type ProxyStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *ProxyStorage {
	return &ProxyStorage{db: db}
}

// If options is nil, all proxies will be returned else only options.PerPage proxies will be returned
// If options.Page is 0, it will be set to 1
// If options.PerPage is 0, it will be set to 10
func (s *ProxyStorage) GetAll(ctx context.Context, options *storage.Options) ([]dto.ProxyDTO, error) {
	var (
		proxies []dto.ProxyDTO
		offset  int
	)

	query := `
        SELECT p.id, p.ip, p.port, p.protocol, p.response_time, p.created_at, p.updated_at,
		   s.id as "status.id", s.name as "status.name",
		   s.created_at as "status.created_at", s.updated_at as "status.updated_at",
		   c.id as "country.id", c.name as "country.name", c.code as "country.code",
		   c.created_at as "country.created_at", c.updated_at as "country.updated_at"
        FROM proxy p
			JOIN status s ON p.status_id = s.id
			JOIN country c ON p.country_id = c.id
		ORDER BY p.id
    `

	if options == nil {
		options = &storage.Options{}
	} else {
		if options.PerPage == 0 {
			options.PerPage = 10
		}
		if options.Page == 0 {
			options.Page = 1
		}
		offset = (options.Page - 1) * options.PerPage
		query += ` LIMIT $1 OFFSET $2`
	}

	err := s.db.SelectContext(ctx, &proxies, query, options.PerPage, offset)
	if err != nil {
		return nil, err
	}

	return proxies, nil
}
