// Package config собирает конфигурацию сервиса из переменных окружения.
package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                LoggerConfig
	Kafka                 KafkaConfig
	OrderPaidConsumer     ConsumerConfig
	ShipAssembledConsumer ConsumerConfig
	Telegram              TelegramConfig
	Metrics               MetricsConfig
	Tracing               TracingConfig
}

// Load читает .env (если он есть) и разбирает переменные окружения.
func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderPaidConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	shipAssembledConsumerCfg, err := env.NewShipAssembledConsumerConfig()
	if err != nil {
		return err
	}

	telegramCfg, err := env.NewTelegramConfig()
	if err != nil {
		return err
	}

	metricsCfg, err := env.NewMetricsConfig()
	if err != nil {
		return err
	}

	tracingCfg, err := env.NewTracingConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                loggerCfg,
		Kafka:                 kafkaCfg,
		OrderPaidConsumer:     orderPaidConsumerCfg,
		ShipAssembledConsumer: shipAssembledConsumerCfg,
		Telegram:              telegramCfg,
		Metrics:               metricsCfg,
		Tracing:               tracingCfg,
	}

	return nil
}

// AppConfig возвращает загруженную конфигурацию.
func AppConfig() *config {
	return appConfig
}
