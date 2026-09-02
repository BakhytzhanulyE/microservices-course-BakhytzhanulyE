// Package app собирает и запускает сервис notification.
package app

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/config"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/closer"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/metrics"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/tracing"
)

// App — приложение assembly.
type App struct {
	diContainer   *diContainer
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

// Run запускает консьюмер и сервер метрик.
func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 3)

	go func() {
		if err := a.diContainer.OrderConsumerService().RunConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("консьюмер order.paid упал: %w", err)
		}
	}()

	go func() {
		if err := a.diContainer.ShipConsumerService().RunConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("консьюмер ship.assembled упал: %w", err)
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
