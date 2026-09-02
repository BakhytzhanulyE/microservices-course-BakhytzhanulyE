// Package closer собирает функции освобождения ресурсов и вызывает их при остановке приложения.
package closer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
)

// shutdownTimeout — сколько ждём закрытия ресурсов по сигналу.
const shutdownTimeout = 5 * time.Second

// Logger — минимальный интерфейс логгера, который нужен closer'у.
type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

// Closer хранит функции закрытия и гарантирует, что они выполнятся ровно один раз.
type Closer struct {
	mu     sync.Mutex
	once   sync.Once
	done   chan struct{}
	funcs  []func(context.Context) error
	logger Logger
}

// globalCloser — экземпляр, которым пользуются все сервисы.
var globalCloser = NewWithLogger(&logger.NoopLogger{})

// AddNamed регистрирует функцию закрытия с именем — оно попадёт в лог.
func AddNamed(name string, f func(context.Context) error) {
	globalCloser.AddNamed(name, f)
}

// Add регистрирует функции закрытия без имени.
func Add(f ...func(context.Context) error) {
	globalCloser.Add(f...)
}

// CloseAll закрывает всё зарегистрированное в глобальном closer'е.
func CloseAll(ctx context.Context) error {
	return globalCloser.CloseAll(ctx)
}

// SetLogger подменяет логгер глобального closer'а.
func SetLogger(l Logger) {
	globalCloser.SetLogger(l)
}

// Configure заставляет глобальный closer слушать системные сигналы.
func Configure(signals ...os.Signal) {
	go globalCloser.handleSignals(signals...)
}

// New создаёт Closer с глобальным логгером.
func New(signals ...os.Signal) *Closer {
	return NewWithLogger(logger.Instance(), signals...)
}

// NewWithLogger создаёт Closer с указанным логгером. Если переданы сигналы — начинает их слушать.
func NewWithLogger(l Logger, signals ...os.Signal) *Closer {
	c := &Closer{
		done:   make(chan struct{}),
		logger: l,
	}

	if len(signals) > 0 {
		go c.handleSignals(signals...)
	}

	return c
}

// SetLogger подменяет логгер.
func (c *Closer) SetLogger(l Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger = l
}

func (c *Closer) log() Logger {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.logger
}

func (c *Closer) handleSignals(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)

	select {
	case <-ch:
		c.log().Info(context.Background(), "🛑 Получен системный сигнал, начинаем graceful shutdown")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := c.CloseAll(shutdownCtx); err != nil {
			c.log().Error(context.Background(), "❌ Ошибка при закрытии ресурсов", zap.Error(err))
		}

	case <-c.done:
		// CloseAll уже отработал вручную — слушать сигналы больше незачем.
	}
}

// AddNamed оборачивает функцию закрытия логированием имени и длительности.
func (c *Closer) AddNamed(name string, f func(context.Context) error) {
	c.Add(func(ctx context.Context) error {
		start := time.Now()
		c.log().Info(ctx, fmt.Sprintf("🧩 Закрываем %s...", name))

		err := f(ctx)

		duration := time.Since(start)
		if err != nil {
			c.log().Error(ctx, fmt.Sprintf("❌ Ошибка при закрытии %s за %s", name, duration), zap.Error(err))
		} else {
			c.log().Info(ctx, fmt.Sprintf("✅ %s закрыт за %s", name, duration))
		}

		return err
	})
}

// Add регистрирует одну или несколько функций закрытия.
func (c *Closer) Add(f ...func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f...)
}

// CloseAll вызывает все зарегистрированные функции в обратном порядке добавления
// и возвращает первую возникшую ошибку.
func (c *Closer) CloseAll(ctx context.Context) error {
	var result error

	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		if len(funcs) == 0 {
			c.log().Info(ctx, "ℹ️ Закрывать нечего")
			return
		}

		c.log().Info(ctx, "🚦 Начинаем graceful shutdown")

		errCh := make(chan error, len(funcs))
		var wg sync.WaitGroup

		for i := len(funcs) - 1; i >= 0; i-- {
			f := funcs[i]
			wg.Add(1)

			go func(f func(context.Context) error) {
				defer wg.Done()

				defer func() {
					if r := recover(); r != nil {
						errCh <- errors.New("panic recovered in closer")
						c.log().Error(ctx, "⚠️ Паника в функции закрытия", zap.Any("error", r))
					}
				}()

				if err := f(ctx); err != nil {
					errCh <- err
				}
			}(f)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for {
			select {
			case <-ctx.Done():
				c.log().Info(ctx, "⚠️ Контекст отменён во время закрытия", zap.Error(ctx.Err()))
				if result == nil {
					result = ctx.Err()
				}

				return

			case err, ok := <-errCh:
				if !ok {
					c.log().Info(ctx, "✅ Все ресурсы закрыты")
					return
				}

				c.log().Error(ctx, "❌ Ошибка при закрытии", zap.Error(err))
				if result == nil {
					result = err
				}
			}
		}
	})

	return result
}
