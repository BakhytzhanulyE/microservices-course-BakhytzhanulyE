// Package cache описывает интерфейс кеша, за которым прячется конкретная реализация.
package cache

import (
	"context"
	"errors"
	"time"
)

// ErrKeyNotFound возвращается, когда ключа нет в кеше.
var ErrKeyNotFound = errors.New("key not found in cache")

// Client — минимальный набор операций кеша, которого хватает сервисам.
type Client interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	Ping(ctx context.Context) error
	Close() error
}
