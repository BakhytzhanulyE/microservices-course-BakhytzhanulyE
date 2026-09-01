package env

import "github.com/caarlos0/env/v11"

type tracingEnvConfig struct {
	// Пустой endpoint выключает трассировку — удобно для локального запуска без Jaeger.
	Endpoint    string `env:"OTEL_EXPORTER_OTLP_ENDPOINT" envDefault:""`
	ServiceName string `env:"OTEL_SERVICE_NAME"           envDefault:"assembly"`
}

type tracingConfig struct {
	raw tracingEnvConfig
}

// NewTracingConfig читает настройки трассировки.
func NewTracingConfig() (*tracingConfig, error) {
	var raw tracingEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &tracingConfig{raw: raw}, nil
}

// Endpoint — адрес OTLP-коллектора.
func (cfg *tracingConfig) Endpoint() string { return cfg.raw.Endpoint }

// ServiceName — имя сервиса в трейсах.
func (cfg *tracingConfig) ServiceName() string { return cfg.raw.ServiceName }
