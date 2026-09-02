package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/converter"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/model"
	iamV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/iam/v1"
)

// Whoami возвращает владельца access-токена.
func (a *api) Whoami(ctx context.Context, req *iamV1.WhoamiRequest) (*iamV1.WhoamiResponse, error) {
	if req.GetAccessToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "access_token is required")
	}

	session, user, err := a.authService.Whoami(ctx, req.GetAccessToken())
	if err != nil {
		if errors.Is(err, model.ErrInvalidToken) {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
		}

		return nil, status.Error(codes.Internal, "failed to resolve token")
	}

	return &iamV1.WhoamiResponse{
		Session: converter.SessionToProto(session),
		User:    converter.UserToProto(user),
	}, nil
}
