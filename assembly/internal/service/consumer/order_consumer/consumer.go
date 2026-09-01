// Package order_consumer читает события оплаченных заказов и запускает сборку.
package order_consumer //nolint:revive // имя пакета повторяет структуру каталогов курса

import (
	"context"
	"time"

	"go.uber.org/zap"

	kafkaConverter "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/assembly/internal/converter/kafka"
	def "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/assembly/internal/service"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderPaidConsumer   kafka.Consumer
	orderPaidDecoder    kafkaConverter.OrderPaidDecoder
	shipProducerService def.ShipProducerService
	buildDuration       time.Duration
}

// NewService создаёт консьюмер оплаченных заказов.
func NewService(
	orderPaidConsumer kafka.Consumer,
	orderPaidDecoder kafkaConverter.OrderPaidDecoder,
	shipProducerService def.ShipProducerService,
	buildDuration time.Duration,
) *service {
	return &service{
		orderPaidConsumer:   orderPaidConsumer,
		orderPaidDecoder:    orderPaidDecoder,
		shipProducerService: shipProducerService,
		buildDuration:       buildDuration,
	}
}

// RunConsumer читает топик до отмены контекста.
func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "🏭 Сборочный цех подписался на оплаченные заказы")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderPaidHandler)
	if err != nil {
		logger.Error(ctx, "Ошибка чтения топика order.paid", zap.Error(err))
		return err
	}

	return nil
}
