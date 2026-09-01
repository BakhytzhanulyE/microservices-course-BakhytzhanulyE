package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type metricsEnvConfig struct {
	Enabled bool   `env:"METRICS_ENABLED" envDefault:"false"`
	Host    string `env:"METRICS_HOST"    envDefault:"0.0.0.0"`
	Port    string `env:"METRICS_PORT"    envDefault:"9090"`
}

type metricsConfig struct {
	raw metricsEnvConfig
}

// NewMetricsConfig читает настройки сервера метрик.
func NewMetricsConfig() (*metricsConfig, error) {
	var raw metricsEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &metricsConfig{raw: raw}, nil
}

// Address — адрес сервера метрик.
func (cfg *metricsConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

// Enabled — поднимать ли сервер метрик.
func (cfg *metricsConfig) Enabled() bool { return cfg.raw.Enabled }
