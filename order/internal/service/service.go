// Package service объявляет бизнес-логику сервиса заказов.
package service

import (
	"context"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
)

// OrderService — операции над заказами.
type OrderService interface {
	Create(ctx context.Context, params model.CreateOrderParams) (model.Order, error)
	Get(ctx context.Context, uuid string) (model.Order, error)
	Pay(ctx context.Context, params model.PayOrderParams) (string, error)
	Cancel(ctx context.Context, uuid string) error
}

// OrderProducerService — отправка доменных событий заказа в Kafka.
type OrderProducerService interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error
}
