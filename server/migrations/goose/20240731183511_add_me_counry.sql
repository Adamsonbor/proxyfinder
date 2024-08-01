-- +goose Up
-- +goose StatementBegin
INSERT INTO country (code, name) VALUES ('ME', 'Montenegro');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM country WHERE code = 'ME';
-- +goose StatementEnd
