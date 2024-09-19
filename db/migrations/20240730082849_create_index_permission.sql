-- +goose Up
-- +goose StatementBegin

-- role permission
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX idx_role_permissions_role_id_permission_id ON role_permissions(role_id, permission_id);

-- user permission
CREATE INDEX idx_user_permissions_user_id ON user_permissions(user_id);
CREATE INDEX idx_user_permissions_permission_id ON user_permissions(permission_id);
CREATE INDEX idx_user_permissions_user_id_permission_id ON user_permissions(user_id, permission_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- role permission
DROP INDEX IF EXISTS idx_role_permissions_role_id;
DROP INDEX IF EXISTS idx_role_permissions_permission_id;
DROP INDEX IF EXISTS idx_role_permissions_role_id_permission_id;

-- user permission
DROP INDEX IF EXISTS idx_user_permissions_user_id;
DROP INDEX IF EXISTS idx_user_permissions_permission_id;
DROP INDEX IF EXISTS idx_user_permissions_user_id_permission_id;

-- +goose StatementEnd
