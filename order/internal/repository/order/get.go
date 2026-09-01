package order

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
	repoConverter "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/repository/converter"
	repoModel "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/repository/model"
)

// Get возвращает заказ по UUID или model.ErrOrderNotFound.
func (r *repository) Get(ctx context.Context, uuid string) (model.Order, error) {
	query := `SELECT ` + selectColumns + ` FROM orders WHERE uuid = $1`

	var order repoModel.Order

	err := r.pool.QueryRow(ctx, query, uuid).Scan(
		&order.UUID,
		&order.UserUUID,
		&order.PartUUIDs,
		&order.TotalPrice,
		&order.TransactionUUID,
		&order.PaymentMethod,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, model.ErrOrderNotFound
		}

		return model.Order{}, err
	}

	return repoConverter.OrderToModel(order), nil
}
