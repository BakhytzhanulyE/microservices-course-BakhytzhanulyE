// Package ship_producer публикует события сборочного цеха в Kafka.
package ship_producer //nolint:revive // имя пакета повторяет структуру каталогов курса

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/assembly/internal/model"
	def "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/assembly/internal/service"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
	eventsV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/events/v1"
)

var _ def.ShipProducerService = (*service)(nil)

type service struct {
	shipAssembledProducer kafka.Producer
}

// NewService создаёт продюсер событий сборки.
func NewService(shipAssembledProducer kafka.Producer) *service {
	return &service{shipAssembledProducer: shipAssembledProducer}
}

// ProduceShipAssembled публикует событие «корабль собран».
func (s *service) ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	msg := &eventsV1.ShipAssembled{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
		AssembledAt:  timestamppb.New(event.AssembledAt),
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "Не удалось сериализовать ShipAssembled", zap.Error(err))
		return err
	}

	return s.shipAssembledProducer.Send(ctx, []byte(event.OrderUUID), payload)
}
