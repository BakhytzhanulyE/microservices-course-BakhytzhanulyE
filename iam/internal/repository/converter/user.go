// Package converter переводит пользователей между доменной моделью и моделью базы.
package converter

import (
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/model"
	repoModel "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/repository/model"
)

// UserToModel переводит строку базы в доменного пользователя.
func UserToModel(user repoModel.User) model.User {
	return model.User{
		UUID:         user.UUID,
		Login:        user.Login,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
	}
}
