-- +goose Up
-- +goose StatementBegin
INSERT INTO users (username, email, password, status, created_at, updated_at) VALUES
('Administrator', 'admin@example.com', 'WCKpRKzlsdtX4T3dKqa8b8D8lazoXew5W2DyFA74hRfzp9//vbKNsbDW7gjpwtU6', 'active', '2024-12-12 10.00', '');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE username = 'Administrator';
-- +goose StatementEnd
