package closer

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
)

func newTestCloser() *Closer {
	return NewWithLogger(&logger.NoopLogger{})
}

func TestCloseAllRunsEveryFunc(t *testing.T) {
	t.Parallel()

	c := newTestCloser()

	var mu sync.Mutex
	closed := make([]int, 0, 3)

	for i := range 3 {
		c.Add(func(_ context.Context) error {
			mu.Lock()
			defer mu.Unlock()
			closed = append(closed, i)

			return nil
		})
	}

	require.NoError(t, c.CloseAll(context.Background()))

	mu.Lock()
	defer mu.Unlock()
	assert.Len(t, closed, 3)
}

func TestCloseAllIsIdempotent(t *testing.T) {
	t.Parallel()

	c := newTestCloser()

	var calls int
	c.Add(func(_ context.Context) error {
		calls++
		return nil
	})

	require.NoError(t, c.CloseAll(context.Background()))
	require.NoError(t, c.CloseAll(context.Background()))

	assert.Equal(t, 1, calls, "повторный CloseAll не должен закрывать ресурсы второй раз")
}

func TestCloseAllReturnsFirstError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("не закрылось")

	c := newTestCloser()
	c.Add(func(_ context.Context) error { return wantErr })

	require.ErrorIs(t, c.CloseAll(context.Background()), wantErr)
}

func TestCloseAllSurvivesPanic(t *testing.T) {
	t.Parallel()

	c := newTestCloser()
	c.Add(func(_ context.Context) error { panic("всё плохо") })

	// Паника в одной функции закрытия не должна ронять процесс.
	err := c.CloseAll(context.Background())

	require.Error(t, err)
}

func TestCloseAllRespectsContext(t *testing.T) {
	t.Parallel()

	c := newTestCloser()

	release := make(chan struct{})
	defer close(release)

	c.Add(func(_ context.Context) error {
		<-release
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := c.CloseAll(ctx)

	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestCloseAllWithoutFuncs(t *testing.T) {
	t.Parallel()

	require.NoError(t, newTestCloser().CloseAll(context.Background()))
}
