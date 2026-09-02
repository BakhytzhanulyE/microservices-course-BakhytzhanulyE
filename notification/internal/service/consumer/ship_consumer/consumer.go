// Package ship_consumer читает события о собранных кораблях.
package ship_consumer //nolint:revive // имя пакета повторяет структуру каталогов курса

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/converter/kafka"
	def "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/service"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/metrics"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	consumer      kafka.Consumer
	decoder       kafkaConverter.ShipAssembledDecoder
	notifyService def.NotifyService
}

// NewService создаёт консьюмер события ShipAssembled.
func NewService(consumer kafka.Consumer, decoder kafkaConverter.ShipAssembledDecoder, notifyService def.NotifyService) *service {
	return &service{
		consumer:      consumer,
		decoder:       decoder,
		notifyService: notifyService,
	}
}

// RunConsumer читает топик до отмены контекста.
func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "🔔 Подписались на собранные корабли")

	err := s.consumer.Consume(ctx, s.handler)
	if err != nil {
		logger.Error(ctx, "Ошибка чтения топика ship.assembled", zap.Error(err))
		return err
	}

	return nil
}

func (s *service) handler(ctx context.Context, msg kafka.Message) (err error) {
	defer func() {
		metrics.ObserveKafka(msg.Topic, err)
	}()

	event, err := s.decoder.Decode(msg.Value)
	if err != nil {
		return err
	}

	logger.Info(ctx, "🚀 Уведомляем о собранном корабле", zap.String("order_uuid", event.OrderUUID))

	return s.notifyService.ShipAssembled(ctx, event)
}
