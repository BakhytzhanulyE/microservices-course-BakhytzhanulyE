// Package tracing настраивает OpenTelemetry: экспорт трейсов по OTLP/gRPC.
package tracing

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials/insecure"
)

const initTimeout = 10 * time.Second

// Init поднимает трейсер и возвращает функцию остановки.
// Если endpoint пустой, трассировка выключена и возвращается no-op закрытие.
func Init(ctx context.Context, serviceName, endpoint string) (func(context.Context) error, error) {
	if endpoint == "" {
		return func(context.Context) error { return nil }, nil
	}

	initCtx, cancel := context.WithTimeout(ctx, initTimeout)
	defer cancel()

	exporter, err := otlptracegrpc.New(initCtx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithTLSCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	res, err := sdkresource.New(initCtx,
		sdkresource.WithAttributes(
			semconv.ServiceName(serviceName),
			attribute.String("environment", "local"),
		),
	)
	if err != nil {
		return nil, err
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return provider.Shutdown, nil
}

// Tracer возвращает трейсер с именем сервиса.
func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}
