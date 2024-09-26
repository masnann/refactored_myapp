-- +goose Up
-- +goose StatementBegin
INSERT INTO users (username, email, password, status, created_at, updated_at) VALUES
('Administrator', 'admin@example.com', 'WCKpRKzlsdtX4T3dKqa8b8D8lazoXew5W2DyFA74hRfzp9//vbKNsbDW7gjpwtU6', 'active', '2024-12-12 10.00', '');

INSERT INTO roles (name, is_active, created_at, updated_at) VALUES
('SuperAdmin', true, '2024-12-12 10.00', ''),
('Customer', true, '2024-12-12 10.00', '');

INSERT INTO permissions (groups, name, created_at, updated_at) VALUES
('PERMISSION', 'CREATE', '2024-12-12 10.00', '');

INSERT INTO user_role (user_id, role_id) VALUES
(1, 1);

INSERT INTO role_permissions (role_id, permission_id, status, created_at, updated_at) VALUES
(1, 1, true, '2024-12-12 20.00', '');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM user_role WHERE user_id = 1 AND role_id = 1;

DELETE FROM role_permissions WHERE role_id = 1 AND permission_id = 1;

DELETE FROM permissions WHERE groups = 'PERMISSION' AND name = 'CREATE' AND created_at = '2024-12-12 10.00';

DELETE FROM roles WHERE name = 'Admin' AND created_at = '2024-12-12 10.00';

DELETE FROM users WHERE username = 'Administrator' AND email = 'admin@example.com';
-- +goose StatementEnd
