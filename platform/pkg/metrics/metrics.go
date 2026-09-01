// Package metrics поднимает реестр Prometheus и отдаёт метрики по HTTP.
package metrics

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	readHeaderTimeout = 5 * time.Second
	metricsPath       = "/metrics"
)

// Registry — реестр метрик приложения. Свой, а не default: так в него не попадает
// ничего лишнего от библиотек и его удобно очищать в тестах.
var (
	Registry = prometheus.NewRegistry()

	grpcRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "app",
			Subsystem: "grpc",
			Name:      "requests_total",
			Help:      "Количество gRPC-запросов по методам и кодам ответа.",
		},
		[]string{"method", "code"},
	)

	grpcDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "app",
			Subsystem: "grpc",
			Name:      "request_duration_seconds",
			Help:      "Длительность обработки gRPC-запросов.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "app",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Количество HTTP-запросов по маршрутам и статусам.",
		},
		[]string{"method", "path", "status"},
	)

	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "app",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "Длительность обработки HTTP-запросов.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	kafkaMessages = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "app",
			Subsystem: "kafka",
			Name:      "messages_total",
			Help:      "Количество обработанных сообщений Kafka по топикам и результату.",
		},
		[]string{"topic", "result"},
	)
)

func init() {
	Registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		grpcRequests,
		grpcDuration,
		httpRequests,
		httpDuration,
		kafkaMessages,
	)
}

// ObserveHTTP фиксирует результат обработки HTTP-запроса.
func ObserveHTTP(method, path, statusCode string, duration time.Duration) {
	httpRequests.WithLabelValues(method, path, statusCode).Inc()
	httpDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

// ObserveKafka фиксирует результат обработки сообщения Kafka.
func ObserveKafka(topic string, err error) {
	result := "ok"
	if err != nil {
		result = "error"
	}

	kafkaMessages.WithLabelValues(topic, result).Inc()
}

// UnaryServerInterceptor считает количество и длительность gRPC-вызовов.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		grpcRequests.WithLabelValues(info.FullMethod, status.Code(err).String()).Inc()
		grpcDuration.WithLabelValues(info.FullMethod).Observe(time.Since(start).Seconds())

		return resp, err
	}
}

// NewServer собирает HTTP-сервер, отдающий метрики на /metrics.
func NewServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle(metricsPath, promhttp.HandlerFor(Registry, promhttp.HandlerOpts{Registry: Registry}))

	return &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}
}

// Serve запускает сервер метрик и не считает ошибкой штатное закрытие.
func Serve(srv *http.Server) error {
	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
