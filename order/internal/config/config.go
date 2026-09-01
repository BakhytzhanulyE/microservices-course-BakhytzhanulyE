// Package config собирает конфигурацию сервиса из переменных окружения.
package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger            LoggerConfig
	HTTP              HTTPConfig
	Postgres          PostgresConfig
	Kafka             KafkaConfig
	OrderPaidProducer OrderPaidProducerConfig
	Clients           ClientsConfig
	Metrics           MetricsConfig
	Tracing           TracingConfig
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

	httpCfg, err := env.NewHTTPConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderPaidProducerCfg, err := env.NewOrderPaidProducerConfig()
	if err != nil {
		return err
	}

	clientsCfg, err := env.NewClientsConfig()
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
		Logger:            loggerCfg,
		HTTP:              httpCfg,
		Postgres:          postgresCfg,
		Kafka:             kafkaCfg,
		OrderPaidProducer: orderPaidProducerCfg,
		Clients:           clientsCfg,
		Metrics:           metricsCfg,
		Tracing:           tracingCfg,
	}

	return nil
}

// AppConfig возвращает загруженную конфигурацию.
func AppConfig() *config {
	return appConfig
}
