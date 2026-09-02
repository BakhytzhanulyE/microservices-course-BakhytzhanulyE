package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/app"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/config"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/closer"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
)

const (
	configPath      = "deploy/compose/order/.env"
	shutdownTimeout = 10 * time.Second
)

// main возвращает код выхода через run: os.Exit нельзя звать напрямую из тела
// main, потому что он не выполняет отложенные вызовы — а на них держится
// graceful shutdown. Ненулевой код обязателен: без него оркестратор и CI
// считают упавший сервис успешно завершившимся.
func main() {
	os.Exit(run())
}

func run() int {
	if err := config.Load(configPath); err != nil {
		panic(fmt.Errorf("не удалось загрузить конфигурацию: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось создать приложение", zap.Error(err))
		return 1
	}

	if err = a.Run(appCtx); err != nil {
		logger.Error(appCtx, "❌ Ошибка при работе приложения", zap.Error(err))
		return 1
	}

	return 0
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
