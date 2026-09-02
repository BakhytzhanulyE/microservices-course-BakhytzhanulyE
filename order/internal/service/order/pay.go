package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
)

// Pay проводит оплату заказа: проверяет статус, дёргает платёжный сервис,
// сохраняет транзакцию и публикует событие OrderPaid.
func (s *service) Pay(ctx context.Context, params model.PayOrderParams) (string, error) {
	order, err := s.orderRepository.Get(ctx, params.OrderUUID)
	if err != nil {
		return "", err
	}

	// Повторная оплата и оплата отменённого заказа — это конфликт, а не ошибка ввода.
	if order.Status != model.OrderStatusPendingPayment {
		return "", model.ErrOrderNotPendingPayment
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, order.UUID, order.UserUUID, params.PaymentMethod)
	if err != nil {
		logger.Error(ctx, "Оплата не прошла", zap.String("order_uuid", order.UUID), zap.Error(err))
		return "", model.ErrPaymentFailed
	}

	paidStatus := model.OrderStatusPaid

	updated, err := s.orderRepository.Update(ctx, order.UUID, model.UpdateOrderParams{
		Status:          &paidStatus,
		TransactionUUID: &transactionUUID,
		PaymentMethod:   &params.PaymentMethod,
	})
	if err != nil {
		return "", err
	}

	// Событие отправляем после того, как статус сохранён: иначе консьюмеры узнают
	// об оплате раньше, чем она зафиксирована в базе.
	err = s.orderProducerService.ProduceOrderPaid(ctx, model.OrderPaidEvent{
		EventUUID:       uuid.NewString(),
		OrderUUID:       updated.UUID,
		UserUUID:        updated.UserUUID,
		PaymentAmount:   updated.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   params.PaymentMethod,
		PaidAt:          time.Now().UTC(),
	})
	if err != nil {
		// Оплата уже прошла, поэтому клиенту отвечаем успехом, а проблему с Kafka
		// оставляем в логах — иначе он попробует заплатить второй раз.
		logger.Error(ctx, "Не удалось отправить событие OrderPaid",
			zap.String("order_uuid", updated.UUID),
			zap.Error(err),
		)
	}

	return transactionUUID, nil
}
