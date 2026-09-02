package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type kafkaEnvConfig struct {
	Brokers []string `env:"KAFKA_BROKERS,required" envSeparator:","`
}

type kafkaConfig struct {
	raw kafkaEnvConfig
}

// NewKafkaConfig читает список брокеров Kafka.
func NewKafkaConfig() (*kafkaConfig, error) {
	var raw kafkaEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &kafkaConfig{raw: raw}, nil
}

// Brokers — адреса брокеров.
func (cfg *kafkaConfig) Brokers() []string { return cfg.raw.Brokers }

type orderPaidProducerEnvConfig struct {
	TopicName string `env:"ORDER_PAID_TOPIC_NAME,required"`
}

type orderPaidProducerConfig struct {
	raw orderPaidProducerEnvConfig
}

// NewOrderPaidProducerConfig читает настройки продюсера событий OrderPaid.
func NewOrderPaidProducerConfig() (*orderPaidProducerConfig, error) {
	var raw orderPaidProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderPaidProducerConfig{raw: raw}, nil
}

// Topic — топик событий OrderPaid.
func (cfg *orderPaidProducerConfig) Topic() string { return cfg.raw.TopicName }

// Config — настройки sarama для синхронного продюсера.
func (cfg *orderPaidProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	// SyncProducer не запустится без Return.Successes: он ждёт подтверждения записи.
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return config
}
