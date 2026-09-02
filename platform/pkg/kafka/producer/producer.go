// Package producer — обёртка над sarama.SyncProducer.
package producer

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka"
)

// tracerName — имя инструментации, под которым спаны видны в Jaeger.
const tracerName = "platform/kafka/producer"

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
//
// Контекст трассировки уезжает вместе с сообщением: пропагатор кладёт traceparent
// в заголовки Kafka, и консьюмер на той стороне продолжает тот же трейс, а не
// начинает свой. Без этого цепочка order → assembly → notification распадалась
// в Jaeger на три несвязанных трейса.
func (p *producer) Send(ctx context.Context, key, value []byte) error {
	ctx, span := otel.Tracer(tracerName).Start(ctx,
		fmt.Sprintf("%s publish", p.topic),
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.operation", "publish"),
			attribute.String("messaging.destination.name", p.topic),
			attribute.String("messaging.kafka.message.key", string(key)),
		),
	)
	defer span.End()

	carrier := kafka.HeaderCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	partition, offset, err := p.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic:   p.topic,
		Key:     sarama.ByteEncoder(key),
		Value:   sarama.ByteEncoder(value),
		Headers: recordHeaders(carrier),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		p.logger.Error(ctx, "Не удалось отправить сообщение в Kafka", zap.Error(err))

		return err
	}

	span.SetAttributes(
		attribute.Int("messaging.kafka.partition", int(partition)),
		attribute.Int64("messaging.kafka.message.offset", offset),
	)

	p.logger.Info(ctx, "Сообщение отправлено",
		zap.String("topic", p.topic),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
		zap.String("key", string(key)),
	)

	return nil
}

// recordHeaders переводит заголовки-карьер в формат sarama.
func recordHeaders(carrier kafka.HeaderCarrier) []sarama.RecordHeader {
	headers := make([]sarama.RecordHeader, 0, len(carrier))
	for k, v := range carrier {
		headers = append(headers, sarama.RecordHeader{Key: []byte(k), Value: v})
	}

	return headers
}
