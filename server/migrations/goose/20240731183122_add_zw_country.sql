-- +goose Up
-- +goose StatementBegin
INSERT INTO country (name, code) VALUES ('Zimbabwe', 'ZW');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM country WHERE code = 'ZW';
-- +goose StatementEnd
