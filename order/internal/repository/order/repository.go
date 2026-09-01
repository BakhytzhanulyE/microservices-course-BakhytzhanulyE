// Package order — реализация хранилища заказов поверх PostgreSQL.
package order

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

// selectColumns перечислены явно: SELECT * ломается при добавлении колонки.
const selectColumns = `uuid, user_uuid, part_uuids, total_price, transaction_uuid, payment_method, status, created_at, updated_at`

type repository struct {
	pool *pgxpool.Pool
}

// NewRepository создаёт хранилище поверх пула соединений с PostgreSQL.
func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}
