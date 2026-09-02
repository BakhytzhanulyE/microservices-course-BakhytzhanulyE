// Package converter переводит заказы между доменной моделью и моделью базы.
package converter

import (
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
	repoModel "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/repository/model"
)

// OrderToModel переводит строку базы в доменный заказ.
func OrderToModel(order repoModel.Order) model.Order {
	return model.Order{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          model.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}
