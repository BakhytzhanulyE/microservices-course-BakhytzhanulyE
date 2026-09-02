// Package decoder разбирает protobuf-сообщения Kafka в доменные события.
package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/model"
	eventsV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/events/v1"
)

type orderPaidDecoder struct{}

// NewOrderPaidDecoder создаёт декодер события OrderPaid.
func NewOrderPaidDecoder() *orderPaidDecoder {
	return &orderPaidDecoder{}
}

// Decode разбирает сообщение Kafka в доменное событие.
func (*orderPaidDecoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsV1.OrderPaid

	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("не удалось разобрать OrderPaid: %w", err)
	}

	return model.OrderPaidEvent{
		EventUUID:       pb.GetEventUuid(),
		OrderUUID:       pb.GetOrderUuid(),
		UserUUID:        pb.GetUserUuid(),
		PaymentAmount:   pb.GetPaymentAmount(),
		TransactionUUID: pb.GetTransactionUuid(),
		PaymentMethod:   pb.GetPaymentMethod(),
		PaidAt:          pb.GetPaidAt().AsTime(),
	}, nil
}
