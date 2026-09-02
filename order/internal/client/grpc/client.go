// Package grpc объявляет контракты клиентов к соседним сервисам.
package grpc

import (
	"context"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
)

// InventoryClient — клиент каталога деталей.
type InventoryClient interface {
	ListParts(ctx context.Context, uuids []string) ([]model.Part, error)
}

// PaymentClient — клиент платёжного сервиса.
type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error)
}
