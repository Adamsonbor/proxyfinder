-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS status (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS status;
-- +goose StatementEnd
