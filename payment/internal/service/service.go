// Package service объявляет бизнес-логику оплаты.
package service

import (
	"context"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/payment/internal/model"
)

// PaymentService — проведение оплаты заказа.
type PaymentService interface {
	PayOrder(ctx context.Context, payment model.Payment) (string, error)
}
