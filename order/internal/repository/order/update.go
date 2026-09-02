package order

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
	repoConverter "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/repository/converter"
	repoModel "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/repository/model"
)

// Update меняет только переданные поля заказа и возвращает результат.
// COALESCE оставляет старое значение там, где параметр пришёл nil,
// поэтому один запрос обслуживает и оплату, и отмену.
func (r *repository) Update(ctx context.Context, uuid string, params model.UpdateOrderParams) (model.Order, error) {
	query := `
		UPDATE orders
		SET status           = COALESCE($2, status),
		    transaction_uuid = COALESCE($3, transaction_uuid),
		    payment_method   = COALESCE($4, payment_method),
		    updated_at       = $5
		WHERE uuid = $1
		RETURNING ` + selectColumns

	var status *string
	if params.Status != nil {
		status = new(string)
		*status = string(*params.Status)
	}

	var order repoModel.Order

	err := r.pool.QueryRow(ctx, query,
		uuid,
		status,
		params.TransactionUUID,
		params.PaymentMethod,
		time.Now().UTC(),
	).Scan(
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
