// Package order реализует бизнес-логику заказов.
package order

import (
	grpcClient "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/client/grpc"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/repository"
	def "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository      repository.OrderRepository
	inventoryClient      grpcClient.InventoryClient
	paymentClient        grpcClient.PaymentClient
	orderProducerService def.OrderProducerService
}

// NewService собирает сервис заказов из хранилища, клиентов соседних сервисов и продюсера событий.
func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient grpcClient.InventoryClient,
	paymentClient grpcClient.PaymentClient,
	orderProducerService def.OrderProducerService,
) *service {
	return &service{
		orderRepository:      orderRepository,
		inventoryClient:      inventoryClient,
		paymentClient:        paymentClient,
		orderProducerService: orderProducerService,
	}
}
