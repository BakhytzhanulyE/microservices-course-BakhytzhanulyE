package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/model"
	iamV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/iam/v1"
)

// Login выдаёт пару токенов по логину и паролю.
func (a *api) Login(ctx context.Context, req *iamV1.LoginRequest) (*iamV1.LoginResponse, error) {
	if req.GetLogin() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "login and password are required")
	}

	tokens, err := a.authService.Login(ctx, model.LoginParams{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	})
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid login or password")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &iamV1.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
