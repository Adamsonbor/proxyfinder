-- +goose Up
-- +goose StatementBegin
INSERT INTO country (code, name) VALUES ('LA', 'Laos');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM country WHERE code = 'LA';
-- +goose StatementEnd
