-- +goose Up
-- +goose StatementBegin
CREATE TABLE OTP (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    otp_hash VARCHAR(255) NOT NULL,
    created_at VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    is_used BOOLEAN DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE OTP;
-- +goose StatementEnd
