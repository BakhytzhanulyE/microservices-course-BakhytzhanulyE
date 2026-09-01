// Package producer — обёртка над sarama.SyncProducer.
package producer

import (
	"context"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka"
)

// Logger — минимальный интерфейс логгера для продюсера.
type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type producer struct {
	syncProducer sarama.SyncProducer
	topic        string
	logger       Logger
}

// NewProducer создаёт продюсер, пишущий в один топик.
func NewProducer(syncProducer sarama.SyncProducer, topic string, logger Logger) kafka.Producer {
	return &producer{
		syncProducer: syncProducer,
		topic:        topic,
		logger:       logger,
	}
}

// Send отправляет сообщение и ждёт подтверждения от брокера.
func (p *producer) Send(ctx context.Context, key, value []byte) error {
	partition, offset, err := p.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	})
	if err != nil {
		p.logger.Error(ctx, "Не удалось отправить сообщение в Kafka", zap.Error(err))
		return err
	}

	p.logger.Info(ctx, "Сообщение отправлено",
		zap.String("topic", p.topic),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
		zap.String("key", string(key)),
	)

	return nil
}
