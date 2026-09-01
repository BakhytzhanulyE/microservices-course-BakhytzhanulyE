package order

import (
	"context"

	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
)

// Cancel отменяет неоплаченный заказ. Оплаченный отменить нельзя.
func (s *service) Cancel(ctx context.Context, uuid string) error {
	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		return err
	}

	switch order.Status {
	case model.OrderStatusPaid:
		return model.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return model.ErrOrderNotPendingPayment
	case model.OrderStatusPendingPayment:
	}

	cancelledStatus := model.OrderStatusCancelled

	if _, err = s.orderRepository.Update(ctx, uuid, model.UpdateOrderParams{Status: &cancelledStatus}); err != nil {
		return err
	}

	logger.Info(ctx, "🚫 Заказ отменён", zap.String("order_uuid", uuid))

	return nil
}
