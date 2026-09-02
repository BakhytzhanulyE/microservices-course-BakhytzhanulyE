// Package repository объявляет контракт хранилища заказов.
package repository

import (
	"context"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
)

// OrderRepository — хранилище заказов.
type OrderRepository interface {
	Create(ctx context.Context, order model.Order) error
	Get(ctx context.Context, uuid string) (model.Order, error)
	Update(ctx context.Context, uuid string, params model.UpdateOrderParams) (model.Order, error)
}
