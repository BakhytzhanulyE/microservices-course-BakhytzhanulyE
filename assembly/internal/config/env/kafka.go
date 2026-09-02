package env

import (
	"time"

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

type orderPaidConsumerEnvConfig struct {
	Topic         string        `env:"ORDER_PAID_TOPIC_NAME,required"`
	GroupID       string        `env:"ORDER_PAID_CONSUMER_GROUP_ID,required"`
	BuildDuration time.Duration `env:"ASSEMBLY_BUILD_DURATION" envDefault:"10s"`
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

// BuildDuration — сколько «собирается» корабль.
func (cfg *orderPaidConsumerConfig) BuildDuration() time.Duration { return cfg.raw.BuildDuration }

// Config — настройки sarama для консьюмера.
func (cfg *orderPaidConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	// OffsetOldest: сервис, поднятый позже остальных, не потеряет уже случившиеся события.
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}

type shipAssembledProducerEnvConfig struct {
	TopicName string `env:"SHIP_ASSEMBLED_TOPIC_NAME,required"`
}

type shipAssembledProducerConfig struct {
	raw shipAssembledProducerEnvConfig
}

// NewShipAssembledProducerConfig читает настройки продюсера ShipAssembled.
func NewShipAssembledProducerConfig() (*shipAssembledProducerConfig, error) {
	var raw shipAssembledProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &shipAssembledProducerConfig{raw: raw}, nil
}

// Topic — топик событий ShipAssembled.
func (cfg *shipAssembledProducerConfig) Topic() string { return cfg.raw.TopicName }

// Config — настройки sarama для синхронного продюсера.
func (cfg *shipAssembledProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return config
}
