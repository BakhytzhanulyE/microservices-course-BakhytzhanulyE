package config

// LoggerConfig — настройки логгера.
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

// GRPCConfig — адрес, который слушает gRPC-сервер.
type GRPCConfig interface {
	Address() string
}

// MongoConfig — подключение к MongoDB.
type MongoConfig interface {
	URI() string
	DatabaseName() string
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
