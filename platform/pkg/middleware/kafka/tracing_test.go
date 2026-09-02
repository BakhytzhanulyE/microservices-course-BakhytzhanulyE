package kafka_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"

	platformKafka "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka"
	kafkaMiddleware "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/middleware/kafka"
)

// TestTracingContinuesProducerTrace проверяет главное свойство: спан консьюмера
// попадает в тот же трейс, что и спан продюсера, и видит его своим родителем.
// Именно этого не хватало — assembly и notification рисовались отдельными трейсами.
func TestTracingContinuesProducerTrace(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Сторона продюсера: заводим спан и кладём его контекст в заголовки.
	producerCtx, producerSpan := provider.Tracer("test").Start(context.Background(), "publish")
	carrier := platformKafka.HeaderCarrier{}
	otel.GetTextMapPropagator().Inject(producerCtx, carrier)
	producerSpan.End()

	require.Contains(t, carrier, "traceparent", "продюсер обязан положить traceparent в заголовки")

	// Сторона консьюмера: контекст сессии пустой, родитель приходит только из заголовков.
	var handlerCtx context.Context
	handler := kafkaMiddleware.Tracing()(func(ctx context.Context, _ platformKafka.Message) error {
		handlerCtx = ctx
		return nil
	})

	err := handler(context.Background(), platformKafka.Message{
		Topic:   "order.paid",
		Headers: carrier,
	})
	require.NoError(t, err)

	consumerSpan := trace.SpanContextFromContext(handlerCtx)
	require.Equal(t, producerSpan.SpanContext().TraceID(), consumerSpan.TraceID(),
		"консьюмер должен продолжать трейс продюсера, а не начинать свой")

	spans := recorder.Ended()
	require.Len(t, spans, 2)
	require.Equal(t, "order.paid process", spans[1].Name())
	require.Equal(t, trace.SpanKindConsumer, spans[1].SpanKind())
	require.Equal(t, producerSpan.SpanContext().SpanID(), spans[1].Parent().SpanID(),
		"родителем спана консьюмера должен быть спан продюсера")
}

// TestTracingWithoutHeaders — сообщение пришло в обход нашего продюсера:
// трейс должен начаться заново, а не упасть.
func TestTracingWithoutHeaders(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	handler := kafkaMiddleware.Tracing()(func(_ context.Context, _ platformKafka.Message) error {
		return nil
	})

	err := handler(context.Background(), platformKafka.Message{Topic: "ship.assembled"})
	require.NoError(t, err)

	spans := recorder.Ended()
	require.Len(t, spans, 1)
	require.False(t, spans[0].Parent().IsValid(), "родителя быть не должно")
}
