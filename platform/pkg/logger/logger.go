// Package logger — обёртка над zap с глобальным экземпляром и обогащением из контекста.
package logger

import (
	"context"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Key — тип ключа для значений, которые логгер достаёт из контекста.
type Key string

const (
	// TraceIDKey — ключ трейса в контексте.
	TraceIDKey Key = "trace_id"
	// UserIDKey — ключ пользователя в контексте.
	UserIDKey Key = "user_id"
)

var (
	globalLogger *Logger
	initOnce     sync.Once
	dynamicLevel zap.AtomicLevel
)

// Logger — обёртка над *zap.Logger, умеющая дописывать поля из контекста.
type Logger struct {
	zapLogger *zap.Logger
}

// Init создаёт глобальный логгер. Повторные вызовы игнорируются.
func Init(levelStr string, asJSON bool) error {
	initOnce.Do(func() {
		dynamicLevel = zap.NewAtomicLevelAt(parseLevel(levelStr))

		encoderCfg := buildEncoderConfig()

		var encoder zapcore.Encoder
		if asJSON {
			encoder = zapcore.NewJSONEncoder(encoderCfg)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderCfg)
		}

		core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), dynamicLevel)

		globalLogger = &Logger{
			// AddCallerSkip(2): один кадр съедает пакетная функция logger.Info,
			// второй — метод (*Logger).Info. Иначе caller всегда указывал бы на этот файл.
			zapLogger: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2)),
		}
	})

	return nil
}

func buildEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

// SetLevel меняет уровень логирования на лету.
func SetLevel(levelStr string) {
	if dynamicLevel == (zap.AtomicLevel{}) {
		return
	}

	dynamicLevel.SetLevel(parseLevel(levelStr))
}

// SetNopLogger переводит глобальный логгер в режим «ничего не пишем» — удобно в тестах.
func SetNopLogger() {
	globalLogger = &Logger{zapLogger: zap.NewNop()}
}

// Instance возвращает глобальный логгер.
func Instance() *Logger {
	if globalLogger == nil {
		return &Logger{zapLogger: zap.NewNop()}
	}

	return globalLogger
}

// Sync сбрасывает буферы логгера.
func Sync() error {
	if globalLogger != nil {
		return globalLogger.zapLogger.Sync()
	}

	return nil
}

// With возвращает логгер с дополнительными постоянными полями.
func With(fields ...zap.Field) *Logger {
	return &Logger{zapLogger: Instance().zapLogger.With(fields...)}
}

// Debug пишет отладочное сообщение через глобальный логгер.
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	Instance().Debug(ctx, msg, fields...)
}

// Info пишет информационное сообщение через глобальный логгер.
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	Instance().Info(ctx, msg, fields...)
}

// Warn пишет предупреждение через глобальный логгер.
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	Instance().Warn(ctx, msg, fields...)
}

// Error пишет ошибку через глобальный логгер.
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	Instance().Error(ctx, msg, fields...)
}

// Debug пишет отладочное сообщение.
func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.zapLogger.Debug(msg, append(fieldsFromContext(ctx), fields...)...)
}

// Info пишет информационное сообщение.
func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.zapLogger.Info(msg, append(fieldsFromContext(ctx), fields...)...)
}

// Warn пишет предупреждение.
func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.zapLogger.Warn(msg, append(fieldsFromContext(ctx), fields...)...)
}

// Error пишет ошибку.
func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.zapLogger.Error(msg, append(fieldsFromContext(ctx), fields...)...)
}

func parseLevel(levelStr string) zapcore.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func fieldsFromContext(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0, 2)

	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		fields = append(fields, zap.String(string(TraceIDKey), traceID))
	}

	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String(string(UserIDKey), userID))
	}

	return fields
}
