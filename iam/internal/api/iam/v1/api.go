// Package v1 — gRPC-слой IAM.
package v1

import (
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/service"
	iamV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/iam/v1"
)

type api struct {
	iamV1.UnimplementedAuthServiceServer

	authService service.AuthService
}

// NewAPI создаёт gRPC-обработчик аутентификации.
func NewAPI(authService service.AuthService) *api {
	return &api{authService: authService}
}
