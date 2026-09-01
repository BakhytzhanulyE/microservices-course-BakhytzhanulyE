// Package order_consumer читает события оплаченных заказов.
package order_consumer //nolint:revive // имя пакета повторяет структуру каталогов курса

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
	decoder       kafkaConverter.OrderPaidDecoder
	notifyService def.NotifyService
}

// NewService создаёт консьюмер события OrderPaid.
func NewService(consumer kafka.Consumer, decoder kafkaConverter.OrderPaidDecoder, notifyService def.NotifyService) *service {
	return &service{
		consumer:      consumer,
		decoder:       decoder,
		notifyService: notifyService,
	}
}

// RunConsumer читает топик до отмены контекста.
func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "🔔 Подписались на оплаченные заказы")

	err := s.consumer.Consume(ctx, s.handler)
	if err != nil {
		logger.Error(ctx, "Ошибка чтения топика order.paid", zap.Error(err))
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

	logger.Info(ctx, "💳 Уведомляем об оплате заказа", zap.String("order_uuid", event.OrderUUID))

	return s.notifyService.OrderPaid(ctx, event)
}
