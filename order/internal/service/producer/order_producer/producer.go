// Package order_producer публикует доменные события заказа в Kafka.
package order_producer //nolint:revive // имя пакета повторяет структуру каталогов курса

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
	def "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/service"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
	eventsV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/events/v1"
)

var _ def.OrderProducerService = (*service)(nil)

type service struct {
	orderPaidProducer kafka.Producer
}

// NewService создаёт продюсер событий заказа.
func NewService(orderPaidProducer kafka.Producer) *service {
	return &service{orderPaidProducer: orderPaidProducer}
}

// ProduceOrderPaid публикует событие «заказ оплачен».
// Ключом берём UUID заказа: так все события одного заказа попадают в одну партицию
// и читаются консьюмером в исходном порядке.
func (s *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	msg := &eventsV1.OrderPaid{
		EventUuid:       event.EventUUID,
		OrderUuid:       event.OrderUUID,
		UserUuid:        event.UserUUID,
		PaymentAmount:   event.PaymentAmount,
		TransactionUuid: event.TransactionUUID,
		PaymentMethod:   event.PaymentMethod,
		PaidAt:          timestamppb.New(event.PaidAt),
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "Не удалось сериализовать OrderPaid", zap.Error(err))
		return err
	}

	return s.orderPaidProducer.Send(ctx, []byte(event.OrderUUID), payload)
}
