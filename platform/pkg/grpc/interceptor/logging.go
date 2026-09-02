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

		switch {
		case err == nil:
			logger.Info(ctx, "gRPC вызов", fields...)
		case isClientError(status.Code(err)):
			logger.Warn(ctx, "gRPC вызов отклонён", append(fields, zap.Error(err))...)
		default:
			logger.Error(ctx, "gRPC вызов завершился ошибкой", append(fields, zap.Error(err))...)
		}

		return resp, err
	}
}

// isClientError отличает «клиент попросил невозможное» от «сервис сломался».
//
// Повторная регистрация или запрос несуществующего заказа — штатный отказ по
// бизнес-правилу, а не сбой. Писать их уровнем ERROR вредно: алерты по частоте
// ошибок начинают срабатывать на обычном поведении пользователей, и настоящая
// авария тонет в этом шуме.
func isClientError(code codes.Code) bool {
	switch code {
	case codes.InvalidArgument,
		codes.NotFound,
		codes.AlreadyExists,
		codes.PermissionDenied,
		codes.Unauthenticated,
		codes.FailedPrecondition,
		codes.OutOfRange,
		codes.Canceled:
		return true
	default:
		return false
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
