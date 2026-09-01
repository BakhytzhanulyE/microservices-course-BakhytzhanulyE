// Package interceptor содержит серверные gRPC-перехватчики: логирование, метрики, recovery.
package interceptor

import (
	"context"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
)

// Logging пишет в лог каждый unary-вызов: метод, код ответа и длительность.
func Logging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.String("code", status.Code(err).String()),
			zap.Duration("duration", time.Since(start)),
		}

		if err != nil {
			logger.Error(ctx, "gRPC вызов завершился ошибкой", append(fields, zap.Error(err))...)
		} else {
			logger.Info(ctx, "gRPC вызов", fields...)
		}

		return resp, err
	}
}

// Recovery ловит панику в обработчике, чтобы она не уронила весь сервер,
// и превращает её в codes.Internal.
func Recovery() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error(ctx, "Паника в gRPC-обработчике",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
					zap.String("stack", string(debug.Stack())),
				)

				err = status.Error(codes.Internal, "internal error")
			}
		}()

		return handler(ctx, req)
	}
}
