-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS protocol (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS protocol;
-- +goose StatementEnd
