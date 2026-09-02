// Package session — хранилище сессий поверх Redis.
package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/model"
	def "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/repository"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/cache"
)

var _ def.SessionRepository = (*repository)(nil)

// keyPrefix отделяет ключи сессий от прочих данных в той же базе Redis.
const keyPrefix = "session:"

type repository struct {
	cache cache.Client
}

// NewRepository создаёт хранилище сессий.
func NewRepository(client cache.Client) *repository {
	return &repository{cache: client}
}

// Create кладёт сессию в Redis с временем жизни ttl.
func (r *repository) Create(ctx context.Context, session model.Session, ttl time.Duration) error {
	return r.cache.Set(ctx, key(session.UUID), session, ttl)
}

// Get возвращает сессию или model.ErrSessionNotFound, если её уже нет.
func (r *repository) Get(ctx context.Context, uuid string) (model.Session, error) {
	data, err := r.cache.Get(ctx, key(uuid))
	if err != nil {
		if errors.Is(err, cache.ErrKeyNotFound) {
			return model.Session{}, model.ErrSessionNotFound
		}

		return model.Session{}, err
	}

	var session model.Session
	if err = json.Unmarshal(data, &session); err != nil {
		return model.Session{}, err
	}

	return session, nil
}

// Delete удаляет сессию — так работает выход из системы.
func (r *repository) Delete(ctx context.Context, uuid string) error {
	return r.cache.Delete(ctx, key(uuid))
}

func key(uuid string) string {
	return fmt.Sprintf("%s%s", keyPrefix, uuid)
}
