// Package payment реализует бизнес-логику оплаты.
package payment

import (
	def "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/payment/internal/service"
)

var _ def.PaymentService = (*service)(nil)

type service struct{}

// NewService создаёт сервис оплаты.
func NewService() *service {
	return &service{}
}
