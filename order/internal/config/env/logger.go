// Package env разбирает переменные окружения в типизированные конфиги.
package env

import "github.com/caarlos0/env/v11"

type loggerEnvConfig struct {
	Level  string `env:"LOGGER_LEVEL"     envDefault:"info"`
	AsJSON bool   `env:"LOGGER_AS_JSON"   envDefault:"true"`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

// NewLoggerConfig читает настройки логгера.
func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &loggerConfig{raw: raw}, nil
}

// Level — уровень логирования.
func (cfg *loggerConfig) Level() string { return cfg.raw.Level }

// AsJSON — писать ли логи в JSON.
func (cfg *loggerConfig) AsJSON() bool { return cfg.raw.AsJSON }
