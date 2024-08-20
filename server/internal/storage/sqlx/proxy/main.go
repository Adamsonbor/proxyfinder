package proxystorage

import (
	"bytes"
	"context"
	"fmt"
	"proxyfinder/internal/domain/dto"
	"proxyfinder/pkg/filter"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	selectAllQuery                  = "SELECT * FROM proxy"
	selectWithStatusAndCountryQuery = `
		SELECT proxy.id, proxy.ip, proxy.port, proxy.protocol,
			proxy.response_time, proxy.status_id, proxy.country_id,
			proxy.created_at, proxy.updated_at,
			status.id as "status.id", status.name as "status.name",
			status.created_at as "status.created_at", status.updated_at as "status.updated_at",
			country.id as "country.id", country.name as "country.name",
			country.created_at as "country.created_at", country.updated_at as "country.updated_at"
		FROM proxy
			JOIN status ON proxy.status_id = status.id
			JOIN country ON proxy.country_id = country.id
		`
	createQuery = `
		INSERT INTO proxy 
		(ip, port, protocol, response_time, status_id, country_id) VALUES 
		(?, ?, ?, ?, ?, ?) RETURNING id
		`
	updateQuery = `UPDATE proxy SET %s WHERE id = ?`
	deleteQuery = `DELETE FROM proxy WHERE id = ?`
)

type ProxyStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) ProxyStorage {
	return ProxyStorage{db: db}
}

// GetAll get options like (status.name, country.code, proxy.response_time)
func (self ProxyStorage) GetAll(ctx context.Context, options filter.Options) ([]dto.Proxy, error) {
	var buf bytes.Buffer
	buf.WriteString(selectWithStatusAndCountryQuery)

	if options.Is() {
		buf.WriteString(" WHERE ")
	}
	for i, field := range options.Fields() {
		if i > 0 {
			buf.WriteString(" AND ")
		}
		buf.WriteString(fmt.Sprintf("%s %s ?", field.Name, field.Op))
	}
	buf.WriteString(" LIMIT ?")
	buf.WriteString(" OFFSET ?")

	values := options.Values()
	values = append(values, options.Limit(), options.Offset())

	var proxies []dto.Proxy
	err := self.db.SelectContext(ctx, &proxies, buf.String(), values...)
	return proxies, err
}

func (self ProxyStorage) Update(ctx context.Context, options filter.Options) error {
	var (
		id interface{}
		query string = ""
	)
	options.AddField("updated_at", filter.OpEq, time.Now(), "time")
	for i, v := range options.Fields() {
		if v.Name == "id" {
			id = v.Val
		}
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("%s = ?", v.Name)
	}
	query = fmt.Sprintf(updateQuery, query)
	values := append(options.Values(), id)

	result, err := self.db.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}
