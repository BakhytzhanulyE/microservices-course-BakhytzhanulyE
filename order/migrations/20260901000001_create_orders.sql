-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    uuid             UUID PRIMARY KEY,
    user_uuid        UUID             NOT NULL,
    part_uuids       TEXT[]           NOT NULL,
    total_price      DOUBLE PRECISION NOT NULL,
    transaction_uuid TEXT,
    payment_method   TEXT,
    status           TEXT             NOT NULL,
    created_at       TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ
);

-- Заказы почти всегда смотрят по пользователю, поэтому индекс именно по нему.
CREATE INDEX IF NOT EXISTS idx_orders_user_uuid ON orders (user_uuid);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders (status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
