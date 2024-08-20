-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS session (
	id INTEGER PRIMARY KEY,
	user_id INTEGER NOT NULL,
	token TEXT NOT NULL,
	expires_at DATETIME NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (user_id, token),
	FOREIGN KEY (user_id) REFERENCES user (id)
		ON UPDATE CASCADE
		ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS session;
-- +goose StatementEnd
