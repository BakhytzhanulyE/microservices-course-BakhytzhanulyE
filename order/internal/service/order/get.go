package order

import (
	"context"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
)

// Get возвращает заказ по UUID.
func (s *service) Get(ctx context.Context, uuid string) (model.Order, error) {
	return s.orderRepository.Get(ctx, uuid)
}
