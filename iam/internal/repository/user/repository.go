// Package user — хранилище пользователей поверх PostgreSQL.
package user

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/repository"
)

var _ def.UserRepository = (*repository)(nil)

const selectColumns = `uuid, login, email, password_hash, created_at`

type repository struct {
	pool *pgxpool.Pool
}

// NewRepository создаёт хранилище пользователей.
func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}
