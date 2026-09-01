package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type redisEnvConfig struct {
	Host     string `env:"REDIS_HOST,required"`
	Port     string `env:"REDIS_PORT,required"`
	Password string `env:"REDIS_PASSWORD" envDefault:""`
	DB       int    `env:"REDIS_DB"       envDefault:"0"`
}

type redisConfig struct {
	raw redisEnvConfig
}

// NewRedisConfig читает настройки подключения к Redis.
func NewRedisConfig() (*redisConfig, error) {
	var raw redisEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &redisConfig{raw: raw}, nil
}

// Address — адрес вида host:port.
func (cfg *redisConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

// Password — пароль Redis.
func (cfg *redisConfig) Password() string { return cfg.raw.Password }

// DB — номер базы Redis.
func (cfg *redisConfig) DB() int { return cfg.raw.DB }
