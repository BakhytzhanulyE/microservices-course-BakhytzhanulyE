package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/model"
	repoConverter "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/repository/converter"
	repoModel "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/repository/model"
)

// GetByUUID возвращает пользователя по UUID.
func (r *repository) GetByUUID(ctx context.Context, uuid string) (model.User, error) {
	return r.getBy(ctx, `SELECT `+selectColumns+` FROM users WHERE uuid = $1`, uuid)
}

// GetByLogin возвращает пользователя по логину.
func (r *repository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	return r.getBy(ctx, `SELECT `+selectColumns+` FROM users WHERE login = $1`, login)
}

func (r *repository) getBy(ctx context.Context, query string, arg any) (model.User, error) {
	var user repoModel.User

	err := r.pool.QueryRow(ctx, query, arg).Scan(
		&user.UUID,
		&user.Login,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		}

		return model.User{}, err
	}

	return repoConverter.UserToModel(user), nil
}
