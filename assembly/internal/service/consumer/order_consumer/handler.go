package order_consumer

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/assembly/internal/model"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/metrics"
)

// OrderPaidHandler собирает корабль по оплаченному заказу и публикует ShipAssembled.
func (s *service) OrderPaidHandler(ctx context.Context, msg kafka.Message) (err error) {
	defer func() {
		metrics.ObserveKafka(msg.Topic, err)
	}()

	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Не удалось разобрать OrderPaid", zap.Error(err))
		return err
	}

	logger.Info(ctx, "🔧 Начинаем сборку корабля",
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.Float64("payment_amount", event.PaymentAmount),
	)

	if err = s.build(ctx); err != nil {
		return err
	}

	assembled := model.ShipAssembledEvent{
		EventUUID:    uuid.NewString(),
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: int64(s.buildDuration.Seconds()),
		AssembledAt:  time.Now().UTC(),
	}

	if err = s.shipProducerService.ProduceShipAssembled(ctx, assembled); err != nil {
		logger.Error(ctx, "Не удалось опубликовать ShipAssembled", zap.Error(err))
		return err
	}

	logger.Info(ctx, "🚀 Корабль собран", zap.String("order_uuid", event.OrderUUID))

	return nil
}

// build имитирует долгую сборку. Ждём через таймер, а не time.Sleep,
// чтобы остановка сервиса не подвисала на длительность сборки.
func (s *service) build(ctx context.Context) error {
	timer := time.NewTimer(s.buildDuration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
