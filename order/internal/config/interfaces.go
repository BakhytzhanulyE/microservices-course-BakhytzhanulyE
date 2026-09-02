package config

import (
	"time"

	"github.com/IBM/sarama"
)

// LoggerConfig — настройки логгера.
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

// HTTPConfig — адрес и таймауты HTTP-сервера.
type HTTPConfig interface {
	Address() string
	ReadHeaderTimeout() time.Duration
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	IdleTimeout() time.Duration
	ShutdownTimeout() time.Duration
}

// PostgresConfig — подключение к PostgreSQL.
type PostgresConfig interface {
	DSN() string
}

// KafkaConfig — адреса брокеров.
type KafkaConfig interface {
	Brokers() []string
}

// OrderPaidProducerConfig — топик и настройки продюсера OrderPaid.
type OrderPaidProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

// ClientsConfig — адреса соседних сервисов.
type ClientsConfig interface {
	InventoryAddress() string
	PaymentAddress() string
}

// MetricsConfig — адрес HTTP-сервера с метриками.
type MetricsConfig interface {
	Address() string
	Enabled() bool
}

// TracingConfig — куда отправлять трейсы.
type TracingConfig interface {
	Endpoint() string
	ServiceName() string
}
