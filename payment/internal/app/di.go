package app

import (
	"context"

	paymentV1API "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/payment/internal/api/payment/v1"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/payment/internal/service"
	paymentService "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/payment/internal/service/payment"
	paymentV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/payment/v1"
)

// diContainer лениво собирает зависимости сервиса.
type diContainer struct {
	paymentV1API paymentV1.PaymentServiceServer

	paymentService service.PaymentService
}

// NewDiContainer создаёт пустой контейнер зависимостей.
func NewDiContainer() *diContainer {
	return &diContainer{}
}

// PaymentV1API возвращает gRPC-обработчик оплаты.
func (d *diContainer) PaymentV1API(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = paymentV1API.NewAPI(d.PaymentService(ctx))
	}

	return d.paymentV1API
}

// PaymentService возвращает бизнес-логику оплаты.
func (d *diContainer) PaymentService(_ context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService()
	}

	return d.paymentService
}
