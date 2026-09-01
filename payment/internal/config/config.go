// Package config собирает конфигурацию сервиса из переменных окружения.
package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/payment/internal/config/env"
)

var appConfig *config

type config struct {
	Logger  LoggerConfig
	GRPC    GRPCConfig
	Metrics MetricsConfig
	Tracing TracingConfig
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

	grpcCfg, err := env.NewGRPCConfig()
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
		Logger:  loggerCfg,
		GRPC:    grpcCfg,
		Metrics: metricsCfg,
		Tracing: tracingCfg,
	}

	return nil
}

// AppConfig возвращает загруженную конфигурацию.
func AppConfig() *config {
	return appConfig
}
