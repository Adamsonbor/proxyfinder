-- +goose Up
-- +goose StatementBegin
INSERT INTO country (code, name) VALUES ('SS', 'South Sudan');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM country WHERE code = 'SS';
-- +goose StatementEnd
