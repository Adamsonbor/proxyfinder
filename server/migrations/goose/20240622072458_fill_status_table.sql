-- +goose Up
-- +goose StatementBegin
INSERT INTO status (name) VALUES 
('not available'),
('available');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM status WHERE name IN ('not available', 'available');
-- +goose StatementEnd
