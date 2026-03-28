-- +goose Up
CREATE TABLE IF NOT EXISTS revoked_tokens (
    jti TEXT PRIMARY KEY,
    expired_at TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS revoked_tokens;
