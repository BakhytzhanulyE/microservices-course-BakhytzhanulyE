package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/model"
)

// uniqueViolationCode — код ошибки PostgreSQL при нарушении уникального индекса.
const uniqueViolationCode = "23505"

// Create сохраняет нового пользователя.
func (r *repository) Create(ctx context.Context, user model.User) error {
	const query = `
		INSERT INTO users (uuid, login, email, password_hash, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.pool.Exec(ctx, query, user.UUID, user.Login, user.Email, user.PasswordHash, user.CreatedAt)
	if err != nil {
		// Уникальный индекс — единственный надёжный способ поймать гонку двух
		// одновременных регистраций: предварительная проверка её не ловит.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return model.ErrUserAlreadyExists
		}

		return err
	}

	return nil
}
