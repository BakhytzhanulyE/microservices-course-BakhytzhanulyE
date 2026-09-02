// Package kafka содержит middleware для консьюмеров Kafka.
package kafka

import (
	"context"
	"time"

	"go.uber.org/zap"

	platformKafka "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka/consumer"
)

// Logger — минимальный интерфейс логгера для middleware.
type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

// Logging пишет в лог факт получения сообщения и время его обработки.
func Logging(logger Logger) consumer.Middleware {
	return func(next platformKafka.MessageHandler) platformKafka.MessageHandler {
		return func(ctx context.Context, msg platformKafka.Message) error {
			start := time.Now()

			logger.Info(ctx, "📨 Получено сообщение Kafka",
				zap.String("topic", msg.Topic),
				zap.Int32("partition", msg.Partition),
				zap.Int64("offset", msg.Offset),
			)

			err := next(ctx, msg)
			if err != nil {
				logger.Error(ctx, "📛 Сообщение обработано с ошибкой",
					zap.String("topic", msg.Topic),
					zap.Duration("duration", time.Since(start)),
					zap.Error(err),
				)

				return err
			}

			logger.Info(ctx, "📬 Сообщение обработано",
				zap.String("topic", msg.Topic),
				zap.Duration("duration", time.Since(start)),
			)

			return nil
		}
	}
}
