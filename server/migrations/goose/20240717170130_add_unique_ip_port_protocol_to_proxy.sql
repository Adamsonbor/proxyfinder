-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS unique_ip_port;
CREATE UNIQUE INDEX IF NOT EXISTS unique_ip_port_protocol ON proxy (ip, port, protocol);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS unique_ip_port_protocol;
-- +goose StatementEnd
