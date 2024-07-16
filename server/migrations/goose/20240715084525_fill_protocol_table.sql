-- +goose Up
-- +goose StatementBegin
INSERT INTO protocol (name) VALUES 
('http'),
('https'),
('socks4'),
('socks5');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM protocol WHERE name IN ('http', 'https', 'socks4', 'socks5');
-- +goose StatementEnd
