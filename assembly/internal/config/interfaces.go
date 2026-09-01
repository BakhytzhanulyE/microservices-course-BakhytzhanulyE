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

// KafkaConfig — адреса брокеров.
type KafkaConfig interface {
	Brokers() []string
}

// OrderPaidConsumerConfig — настройки консьюмера OrderPaid.
type OrderPaidConsumerConfig interface {
	Topic() string
	GroupID() string
	BuildDuration() time.Duration
	Config() *sarama.Config
}

// ShipAssembledProducerConfig — настройки продюсера ShipAssembled.
type ShipAssembledProducerConfig interface {
	Topic() string
	Config() *sarama.Config
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
