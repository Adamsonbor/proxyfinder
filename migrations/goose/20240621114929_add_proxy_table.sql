-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS proxy (
	id INTEGER PRIMARY KEY,
	ip TEXT NOT NULL,
	port INTEGER,
	protocol TEXT,
	status_id INTEGER,
	country_id INTEGER,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS proxy;
-- +goose StatementEnd
