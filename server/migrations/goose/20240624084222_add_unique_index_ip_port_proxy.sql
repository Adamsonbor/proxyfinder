-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS unique_ip_port ON proxy (ip, port);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS unique_ip_port;
-- +goose StatementEnd
