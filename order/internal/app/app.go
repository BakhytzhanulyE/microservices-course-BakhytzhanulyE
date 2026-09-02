// Package app собирает и запускает сервис order.
package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/config"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/migrations"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/closer"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/metrics"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/migrator"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/tracing"
)

// App — приложение order со всеми его серверами.
type App struct {
	diContainer   *diContainer
	httpServer    *http.Server
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
		logger.Info(ctx, fmt.Sprintf("🚀 HTTP OrderService слушает на %s", config.AppConfig().HTTP.Address()))

		err := a.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("HTTP-сервер упал: %w", err)
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
		a.initMigrations,
		a.initHTTPServer,
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

// initMigrations накатывает схему до того, как сервис начнёт принимать запросы.
func (a *App) initMigrations(ctx context.Context) error {
	if err := migrator.Up(ctx, config.AppConfig().Postgres.DSN(), migrations.FS, "."); err != nil {
		return fmt.Errorf("не удалось применить миграции: %w", err)
	}

	logger.Info(ctx, "📜 Миграции применены")

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	httpCfg := config.AppConfig().HTTP

	a.httpServer = &http.Server{
		Addr: httpCfg.Address(),
		// otelhttp заводит корневой спан запроса — от него дальше наследуются
		// спаны вызовов inventory и payment.
		Handler:           otelhttp.NewHandler(a.diContainer.OrderV1API(ctx).Router(), "order-http"),
		ReadHeaderTimeout: httpCfg.ReadHeaderTimeout(),
		ReadTimeout:       httpCfg.ReadTimeout(),
		WriteTimeout:      httpCfg.WriteTimeout(),
		IdleTimeout:       httpCfg.IdleTimeout(),
	}

	// Shutdown, а не Close: активные запросы должны успеть дойти до ответа.
	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(ctx, httpCfg.ShutdownTimeout())
		defer cancel()

		return a.httpServer.Shutdown(shutdownCtx)
	})

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
