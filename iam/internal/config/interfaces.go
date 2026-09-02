package config

import "time"

// LoggerConfig — настройки логгера.
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

// GRPCConfig — адрес, который слушает gRPC-сервер.
type GRPCConfig interface {
	Address() string
}

// PostgresConfig — подключение к PostgreSQL.
type PostgresConfig interface {
	DSN() string
}

// RedisConfig — подключение к Redis.
type RedisConfig interface {
	Address() string
	Password() string
	DB() int
}

// JWTConfig — ключ подписи и время жизни токенов.
type JWTConfig interface {
	SecretKey() string
	AccessTTL() time.Duration
	RefreshTTL() time.Duration
}

// MetricsConfig — адрес HTTP-сервера с метриками.
type MetricsConfig interface {
	Address() string
	Enabled() bool
}

// TracingConfig — куда отправлять трейсы.
type TracingConfig interface {
	Endpoint() string
	ServiceName() string
}
