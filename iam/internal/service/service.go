// Package service объявляет бизнес-логику IAM.
package service

import (
	"context"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/model"
)

// AuthService — регистрация, вход и проверка токена.
type AuthService interface {
	Register(ctx context.Context, params model.RegisterParams) (string, error)
	Login(ctx context.Context, params model.LoginParams) (model.TokenPair, error)
	Whoami(ctx context.Context, accessToken string) (model.Session, model.User, error)
}
