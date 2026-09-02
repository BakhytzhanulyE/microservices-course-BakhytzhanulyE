package kafka

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	platformKafka "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka/consumer"
)

// tracerName — имя инструментации, под которым спаны видны в Jaeger.
const tracerName = "platform/kafka/consumer"

// Tracing продолжает трейс, начатый продюсером.
//
// Консьюмер получает от sarama session.Context() — контекст сессии группы, в нём
// нет никакого спана. Родителя достаём из заголовков сообщения (traceparent их
// туда положил продюсер) и передаём обработчику уже контекст со спаном, поэтому
// вся работа сервиса ложится внутрь общего трейса.
//
// Ставить эту middleware нужно первой в цепочке: тогда всё остальное — логирование,
// декодирование, бизнес-логика — попадает внутрь спана.
func Tracing() consumer.Middleware {
	return func(next platformKafka.MessageHandler) platformKafka.MessageHandler {
		return func(ctx context.Context, msg platformKafka.Message) error {
			// Extract возвращает контекст с удалённым родителем; если заголовка нет
			// (например, сообщение прислали в обход нашего продюсера), спан просто
			// станет корнем нового трейса.
			ctx = otel.GetTextMapPropagator().Extract(ctx, platformKafka.HeaderCarrier(msg.Headers))

			ctx, span := otel.Tracer(tracerName).Start(ctx,
				fmt.Sprintf("%s process", msg.Topic),
				trace.WithSpanKind(trace.SpanKindConsumer),
				trace.WithAttributes(
					attribute.String("messaging.system", "kafka"),
					attribute.String("messaging.operation", "process"),
					attribute.String("messaging.destination.name", msg.Topic),
					attribute.String("messaging.kafka.message.key", string(msg.Key)),
					attribute.Int("messaging.kafka.partition", int(msg.Partition)),
					attribute.Int64("messaging.kafka.message.offset", msg.Offset),
				),
			)
			defer span.End()

			if err := next(ctx, msg); err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())

				return err
			}

			return nil
		}
	}
}
