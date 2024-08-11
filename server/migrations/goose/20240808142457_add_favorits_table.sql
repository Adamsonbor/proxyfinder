-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS favorits (
	id INTEGER PRIMARY KEY,
	user_id INTEGER NOT NULL,
	proxy_id INTEGER NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(user_id, proxy_id),
	FOREIGN KEY (user_id) REFERENCES user(id),
	FOREIGN KEY (proxy_id) REFERENCES proxy(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS favorits;
-- +goose StatementEnd
