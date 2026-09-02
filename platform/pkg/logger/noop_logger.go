package logger

import (
	"context"

	"go.uber.org/zap"
)

// NoopLogger — логгер-заглушка, ничего не пишет. Нужен до вызова Init.
type NoopLogger struct{}

// Info ничего не делает.
func (*NoopLogger) Info(_ context.Context, _ string, _ ...zap.Field) {}

// Error ничего не делает.
func (*NoopLogger) Error(_ context.Context, _ string, _ ...zap.Field) {}
