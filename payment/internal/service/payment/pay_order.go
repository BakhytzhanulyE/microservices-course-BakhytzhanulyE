package payment

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/payment/internal/model"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
)

// PayOrder «проводит» оплату и возвращает UUID транзакции.
// Настоящего платёжного шлюза здесь нет: сервис учебный и просто фиксирует факт оплаты.
func (s *service) PayOrder(ctx context.Context, payment model.Payment) (string, error) {
	if !payment.PaymentMethod.Valid() {
		return "", model.ErrUnsupportedPaymentMethod
	}

	transactionUUID := uuid.NewString()

	logger.Info(ctx, "💳 Оплата прошла успешно",
		zap.String("order_uuid", payment.OrderUUID),
		zap.String("user_uuid", payment.UserUUID),
		zap.String("payment_method", payment.PaymentMethod.String()),
		zap.String("transaction_uuid", transactionUUID),
	)

	return transactionUUID, nil
}
