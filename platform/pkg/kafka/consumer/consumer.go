// Package consumer — обёртка над sarama.ConsumerGroup с поддержкой middleware.
package consumer

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka"
)

// Logger — минимальный интерфейс логгера для консьюмера.
type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type consumer struct {
	group       sarama.ConsumerGroup
	topics      []string
	logger      Logger
	middlewares []Middleware
}

// NewConsumer создаёт консьюмер поверх готовой consumer group.
func NewConsumer(group sarama.ConsumerGroup, topics []string, logger Logger, middlewares ...Middleware) kafka.Consumer {
	return &consumer{
		group:       group,
		topics:      topics,
		logger:      logger,
		middlewares: middlewares,
	}
}

// Consume читает топики до отмены контекста. Ребаланс группы — это выход из
// group.Consume без ошибки, поэтому вызов обёрнут в цикл.
func (c *consumer) Consume(ctx context.Context, handler kafka.MessageHandler) error {
	groupHandler := NewGroupHandler(handler, c.logger, c.middlewares...)

	for {
		if err := c.group.Consume(ctx, c.topics, groupHandler); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}

			c.logger.Error(ctx, "Ошибка чтения из Kafka", zap.Error(err))

			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		c.logger.Info(ctx, "Ребалансировка consumer group")
	}
}
