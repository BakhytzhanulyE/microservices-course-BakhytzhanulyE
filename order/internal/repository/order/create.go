package order

import (
	"context"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
)

// Create сохраняет новый заказ.
func (r *repository) Create(ctx context.Context, order model.Order) error {
	const query = `
		INSERT INTO orders (uuid, user_uuid, part_uuids, total_price, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.pool.Exec(ctx, query,
		order.UUID,
		order.UserUUID,
		order.PartUUIDs,
		order.TotalPrice,
		string(order.Status),
		order.CreatedAt,
	)

	return err
}
