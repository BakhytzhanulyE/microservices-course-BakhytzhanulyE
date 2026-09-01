// Package app собирает и запускает сервис payment.
package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/payment/internal/config"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/closer"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/grpc/health"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/grpc/interceptor"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/metrics"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/tracing"
	paymentV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/payment/v1"
)

// App — приложение payment.
type App struct {
	diContainer   *diContainer
	grpcServer    *grpc.Server
	listener      net.Listener
	metricsServer *http.Server
}

// New создаёт приложение и инициализирует зависимости.
func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

// Run запускает серверы и ждёт либо ошибки, либо отмены контекста.
func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 2)

	go func() {
		if err := a.runGRPCServer(ctx); err != nil {
			errCh <- fmt.Errorf("gRPC-сервер упал: %w", err)
		}
	}()

	if a.metricsServer != nil {
		go func() {
			if err := metrics.Serve(a.metricsServer); err != nil {
				errCh <- fmt.Errorf("сервер метрик упал: %w", err)
			}
		}()
	}

	select {
	case <-ctx.Done():
		logger.Info(ctx, "Получен сигнал остановки")
		return nil

	case err := <-errCh:
		logger.Error(ctx, "Компонент упал, останавливаем сервис", zap.Error(err))
		cancel()

		return err
	}
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initTracing,
		a.initListener,
		a.initGRPCServer,
		a.initMetricsServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJSON(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Instance())
	return nil
}

func (a *App) initTracing(ctx context.Context) error {
	shutdown, err := tracing.Init(ctx,
		config.AppConfig().Tracing.ServiceName(),
		config.AppConfig().Tracing.Endpoint(),
	)
	if err != nil {
		return err
	}

	closer.AddNamed("Tracing", shutdown)

	return nil
}

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().GRPC.Address())
	if err != nil {
		return err
	}

	closer.AddNamed("TCP listener", func(_ context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})

	a.listener = listener

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			interceptor.Recovery(),
			interceptor.Logging(),
			metrics.UnaryServerInterceptor(),
		),
	)

	closer.AddNamed("gRPC server", func(_ context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)
	health.RegisterService(a.grpcServer)

	paymentV1.RegisterPaymentServiceServer(a.grpcServer, a.diContainer.PaymentV1API(ctx))

	return nil
}

func (a *App) initMetricsServer(_ context.Context) error {
	if !config.AppConfig().Metrics.Enabled() {
		return nil
	}

	a.metricsServer = metrics.NewServer(config.AppConfig().Metrics.Address())

	closer.AddNamed("Metrics server", func(ctx context.Context) error {
		return a.metricsServer.Shutdown(ctx)
	})

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC PaymentService слушает на %s", config.AppConfig().GRPC.Address()))

	return a.grpcServer.Serve(a.listener)
}
