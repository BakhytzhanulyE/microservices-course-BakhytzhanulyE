package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type httpEnvConfig struct {
	Host              string        `env:"HTTP_HOST,required"`
	Port              string        `env:"HTTP_PORT,required"`
	ReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT" envDefault:"5s"`
	ReadTimeout       time.Duration `env:"HTTP_READ_TIMEOUT"        envDefault:"10s"`
	WriteTimeout      time.Duration `env:"HTTP_WRITE_TIMEOUT"       envDefault:"10s"`
	IdleTimeout       time.Duration `env:"HTTP_IDLE_TIMEOUT"        envDefault:"60s"`
	ShutdownTimeout   time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT"    envDefault:"10s"`
}

type httpConfig struct {
	raw httpEnvConfig
}

// NewHTTPConfig читает настройки HTTP-сервера.
func NewHTTPConfig() (*httpConfig, error) {
	var raw httpEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &httpConfig{raw: raw}, nil
}

// Address — адрес вида host:port.
func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

// ReadHeaderTimeout — таймаут чтения заголовков.
func (cfg *httpConfig) ReadHeaderTimeout() time.Duration { return cfg.raw.ReadHeaderTimeout }

// ReadTimeout — таймаут чтения запроса.
func (cfg *httpConfig) ReadTimeout() time.Duration { return cfg.raw.ReadTimeout }

// WriteTimeout — таймаут записи ответа.
func (cfg *httpConfig) WriteTimeout() time.Duration { return cfg.raw.WriteTimeout }

// IdleTimeout — таймаут простоя keep-alive соединения.
func (cfg *httpConfig) IdleTimeout() time.Duration { return cfg.raw.IdleTimeout }

// ShutdownTimeout — сколько ждём завершения активных запросов при остановке.
func (cfg *httpConfig) ShutdownTimeout() time.Duration { return cfg.raw.ShutdownTimeout }
