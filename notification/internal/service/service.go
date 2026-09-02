// Package service объявляет бизнес-логику уведомлений.
package service

import (
	"context"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/model"
)

// ConsumerService читает события из Kafka до отмены контекста.
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

// NotifyService отправляет уведомления пользователю.
type NotifyService interface {
	OrderPaid(ctx context.Context, event model.OrderPaidEvent) error
	ShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}
