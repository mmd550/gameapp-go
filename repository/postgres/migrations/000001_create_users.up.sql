CREATE TABLE IF NOT EXISTS users (
    id          BIGSERIAL    PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    name        TEXT         NOT NULL,
    phone_number TEXT        NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);
