-- +goose Up
-- +goose StatementBegin
INSERT INTO country (name, code) VALUES ('Kosovo', 'XK');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM country WHERE code = 'XK';
-- +goose StatementEnd
