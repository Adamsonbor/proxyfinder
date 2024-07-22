-- +goose Up
-- +goose StatementBegin
INSERT INTO status (name) VALUES 
('Available'),
('Unavailable');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM status WHERE name IN ('not available', 'available');
-- +goose StatementEnd
