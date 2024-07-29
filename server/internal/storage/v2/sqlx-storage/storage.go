package sqlxstorage

import (
	"context"
	"proxyfinder/internal/storage/v2/dto"

	"github.com/jmoiron/sqlx"
)

type ProxyStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *ProxyStorage {
	return &ProxyStorage{db: db}
}

func (s *ProxyStorage) GetAll(ctx context.Context, page int, perPage int) ([]dto.ProxyDTO, error) {
    var proxies []dto.ProxyDTO

	if page == 0 {
		page = 1
	}

	offset := (page - 1) * perPage

    query := `
        SELECT p.id, p.ip, p.port, p.protocol, p.response_time,
		   s.id as "status.id", s.name as "status.name",
		   s.created_at as "status.created_at", s.updated_at as "status.updated_at",
		   c.id as "country.id", c.name as "country.name", c.code as "country.code",
		   c.created_at as "country.created_at", c.updated_at as "country.updated_at"
        FROM proxy p
			JOIN status s ON p.status_id = s.id
			JOIN country c ON p.country_id = c.id
		ORDER BY p.id
		LIMIT $1 OFFSET $2
    `

    err := s.db.SelectContext(ctx, &proxies, query, perPage, offset)
    if err != nil {
        return nil, err
    }

    return proxies, nil
}
