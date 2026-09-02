package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/model"
	iamV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/iam/v1"
)

// Register регистрирует нового пользователя.
func (a *api) Register(ctx context.Context, req *iamV1.RegisterRequest) (*iamV1.RegisterResponse, error) {
	if req.GetLogin() == "" || req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "login, email and password are required")
	}

	userUUID, err := a.authService.Register(ctx, model.RegisterParams{
		Login:    req.GetLogin(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUserAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		case errors.Is(err, model.ErrWeakPassword):
			return nil, status.Error(codes.InvalidArgument, "password is too short")
		default:
			return nil, status.Error(codes.Internal, "failed to register user")
		}
	}

	return &iamV1.RegisterResponse{UserUuid: userUUID}, nil
}
