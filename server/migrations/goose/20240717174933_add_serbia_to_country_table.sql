-- +goose Up
-- +goose StatementBegin
INSERT INTO country (name, code) VALUES ('Serbia', 'RS');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM country WHERE code = 'RS';
-- +goose StatementEnd
