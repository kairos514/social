-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY,
    email CITEXT UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password BYTEA NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
--SELECT 'down SQL query';
-- +goose StatementEnd
