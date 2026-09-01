package config

import "github.com/IBM/sarama"

// LoggerConfig — настройки логгера.
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

// KafkaConfig — адреса брокеров.
type KafkaConfig interface {
	Brokers() []string
}

// ConsumerConfig — настройки одного консьюмера.
type ConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

// TelegramConfig — настройки Telegram-бота.
type TelegramConfig interface {
	Token() string
	ChatID() string
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
