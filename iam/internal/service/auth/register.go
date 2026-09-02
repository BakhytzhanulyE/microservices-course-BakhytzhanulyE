package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/model"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
)

// Register создаёт пользователя и возвращает его UUID.
func (s *service) Register(ctx context.Context, params model.RegisterParams) (string, error) {
	if len(params.Password) < minPasswordLength {
		return "", model.ErrWeakPassword
	}

	// bcrypt солит хеш сам, поэтому отдельное поле для соли не нужно.
	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := model.User{
		UUID:         uuid.NewString(),
		Login:        params.Login,
		Email:        params.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now().UTC(),
	}

	if err = s.userRepository.Create(ctx, user); err != nil {
		return "", err
	}

	logger.Info(ctx, "👤 Пользователь зарегистрирован", zap.String("user_uuid", user.UUID))

	return user.UUID, nil
}
