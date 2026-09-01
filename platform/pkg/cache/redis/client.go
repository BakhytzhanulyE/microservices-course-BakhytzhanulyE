// Package redis — реализация cache.Client поверх go-redis.
package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/cache"
)

type client struct {
	rdb *redis.Client
}

// NewClient создаёт клиент Redis по адресу addr.
func NewClient(addr, password string, db int) cache.Client {
	return &client{
		rdb: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

// Set кладёт значение в кеш. Всё, кроме строк и байтов, сериализуется в JSON.
func (c *client) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	var payload any

	switch v := value.(type) {
	case string, []byte:
		payload = v
	default:
		encoded, err := json.Marshal(v)
		if err != nil {
			return err
		}

		payload = encoded
	}

	return c.rdb.Set(ctx, key, payload, ttl).Err()
}

// Get возвращает значение по ключу или cache.ErrKeyNotFound.
func (c *client) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, cache.ErrKeyNotFound
		}

		return nil, err
	}

	return res, nil
}

// Delete удаляет ключи.
func (c *client) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	return c.rdb.Del(ctx, keys...).Err()
}

// Exists говорит, есть ли ключ в кеше.
func (c *client) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Ping проверяет доступность Redis.
func (c *client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// Close закрывает соединение.
func (c *client) Close() error {
	return c.rdb.Close()
}
