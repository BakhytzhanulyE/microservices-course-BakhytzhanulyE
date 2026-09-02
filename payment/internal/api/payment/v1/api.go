// Package v1 — gRPC-слой сервиса оплаты.
package v1

import (
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/payment/internal/service"
	paymentV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

// NewAPI создаёт gRPC-обработчик оплаты.
func NewAPI(paymentService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
