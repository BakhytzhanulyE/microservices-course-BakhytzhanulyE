-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    uuid          UUID PRIMARY KEY,
    login         TEXT        NOT NULL,
    email         TEXT        NOT NULL,
    password_hash TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Логин и почта должны быть уникальны: без этих индексов гонка двух регистраций
-- создала бы двух пользователей с одним логином.
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_login ON users (login);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
