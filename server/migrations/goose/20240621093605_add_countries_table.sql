-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS country (
	id INTEGER PRIMARY KEY,
	name TEXT,
	code TEXT,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS country;
-- +goose StatementEnd
