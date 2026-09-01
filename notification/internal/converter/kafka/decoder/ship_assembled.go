package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/model"
	eventsV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/events/v1"
)

type shipAssembledDecoder struct{}

// NewShipAssembledDecoder создаёт декодер события ShipAssembled.
func NewShipAssembledDecoder() *shipAssembledDecoder {
	return &shipAssembledDecoder{}
}

// Decode разбирает сообщение Kafka в доменное событие.
func (*shipAssembledDecoder) Decode(data []byte) (model.ShipAssembledEvent, error) {
	var pb eventsV1.ShipAssembled

	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("не удалось разобрать ShipAssembled: %w", err)
	}

	return model.ShipAssembledEvent{
		EventUUID:    pb.GetEventUuid(),
		OrderUUID:    pb.GetOrderUuid(),
		UserUUID:     pb.GetUserUuid(),
		BuildTimeSec: pb.GetBuildTimeSec(),
		AssembledAt:  pb.GetAssembledAt().AsTime(),
	}, nil
}
