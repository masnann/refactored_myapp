-- +goose Up
-- +goose StatementBegin
CREATE TABLE partner (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    secret_key VARCHAR(255) NOT NULL,
    public_key VARCHAR(500),
    domain VARCHAR(255),
    created_at VARCHAR(255),
    updated_at VARCHAR(255),
    deleted_at VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE partner;
-- +goose StatementEnd
