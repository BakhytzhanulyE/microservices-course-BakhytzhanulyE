// Package repository объявляет контракты хранилищ IAM.
package repository

import (
	"context"
	"time"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/model"
)

// UserRepository — хранилище пользователей.
type UserRepository interface {
	Create(ctx context.Context, user model.User) error
	GetByUUID(ctx context.Context, uuid string) (model.User, error)
	GetByLogin(ctx context.Context, login string) (model.User, error)
}

// SessionRepository — хранилище сессий. Живёт в Redis: сессия временная,
// переживать перезапуск сервиса ей не нужно, а TTL Redis делает сам.
type SessionRepository interface {
	Create(ctx context.Context, session model.Session, ttl time.Duration) error
	Get(ctx context.Context, uuid string) (model.Session, error)
	Delete(ctx context.Context, uuid string) error
}
