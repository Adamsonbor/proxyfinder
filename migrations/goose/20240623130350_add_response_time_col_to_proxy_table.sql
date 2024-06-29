-- +goose Up
-- +goose StatementBegin
ALTER TABLE proxy ADD response_time INTEGER DEFAULT 100000;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE proxy DROP COLUMN response_time;
-- +goose StatementEnd
