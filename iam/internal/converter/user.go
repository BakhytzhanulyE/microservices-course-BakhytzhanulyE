// Package converter переводит сущности IAM в protobuf.
package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/model"
	iamV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/iam/v1"
)

// UserToProto переводит пользователя в protobuf. Хеш пароля наружу не уезжает.
func UserToProto(user model.User) *iamV1.User {
	return &iamV1.User{
		Uuid:      user.UUID,
		Login:     user.Login,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}

// SessionToProto переводит сессию в protobuf.
func SessionToProto(session model.Session) *iamV1.Session {
	return &iamV1.Session{
		Uuid:      session.UUID,
		UserUuid:  session.UserUUID,
		ExpiresAt: timestamppb.New(session.ExpiresAt),
	}
}
