// Package auth реализует регистрацию, вход и проверку токенов.
package auth

import (
	"time"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/repository"
	def "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/service"
)

var _ def.AuthService = (*service)(nil)

// minPasswordLength — минимальная длина пароля.
const minPasswordLength = 8

type service struct {
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository

	secretKey  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewService собирает сервис аутентификации.
func NewService(
	userRepository repository.UserRepository,
	sessionRepository repository.SessionRepository,
	secretKey string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *service {
	return &service{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		secretKey:         secretKey,
		accessTTL:         accessTTL,
		refreshTTL:        refreshTTL,
	}
}
