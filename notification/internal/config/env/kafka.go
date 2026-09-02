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

// consumerConfig — общие настройки sarama для консьюмеров сервиса.
func consumerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}

type orderPaidConsumerEnvConfig struct {
	Topic   string `env:"ORDER_PAID_TOPIC_NAME,required"`
	GroupID string `env:"ORDER_PAID_CONSUMER_GROUP_ID,required"`
}

type orderPaidConsumerConfig struct {
	raw orderPaidConsumerEnvConfig
}

// NewOrderPaidConsumerConfig читает настройки консьюмера OrderPaid.
func NewOrderPaidConsumerConfig() (*orderPaidConsumerConfig, error) {
	var raw orderPaidConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderPaidConsumerConfig{raw: raw}, nil
}

// Topic — топик, из которого читаем.
func (cfg *orderPaidConsumerConfig) Topic() string { return cfg.raw.Topic }

// GroupID — идентификатор consumer group.
func (cfg *orderPaidConsumerConfig) GroupID() string { return cfg.raw.GroupID }

// Config — настройки sarama.
func (cfg *orderPaidConsumerConfig) Config() *sarama.Config { return consumerConfig() }

type shipAssembledConsumerEnvConfig struct {
	Topic   string `env:"SHIP_ASSEMBLED_TOPIC_NAME,required"`
	GroupID string `env:"SHIP_ASSEMBLED_CONSUMER_GROUP_ID,required"`
}

type shipAssembledConsumerConfig struct {
	raw shipAssembledConsumerEnvConfig
}

// NewShipAssembledConsumerConfig читает настройки консьюмера ShipAssembled.
func NewShipAssembledConsumerConfig() (*shipAssembledConsumerConfig, error) {
	var raw shipAssembledConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &shipAssembledConsumerConfig{raw: raw}, nil
}

// Topic — топик, из которого читаем.
func (cfg *shipAssembledConsumerConfig) Topic() string { return cfg.raw.Topic }

// GroupID — идентификатор consumer group.
func (cfg *shipAssembledConsumerConfig) GroupID() string { return cfg.raw.GroupID }

// Config — настройки sarama.
func (cfg *shipAssembledConsumerConfig) Config() *sarama.Config { return consumerConfig() }
