package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/model"
)

// Login проверяет пароль, заводит сессию и выдаёт пару токенов.
func (s *service) Login(ctx context.Context, params model.LoginParams) (model.TokenPair, error) {
	user, err := s.userRepository.GetByLogin(ctx, params.Login)
	if err != nil {
		// Не различаем «нет пользователя» и «неверный пароль»: иначе по ответу
		// можно перебором узнать, какие логины существуют.
		if errors.Is(err, model.ErrUserNotFound) {
			return model.TokenPair{}, model.ErrInvalidCredentials
		}

		return model.TokenPair{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password)); err != nil {
		return model.TokenPair{}, model.ErrInvalidCredentials
	}

	session := model.Session{
		UUID:      uuid.NewString(),
		UserUUID:  user.UUID,
		ExpiresAt: time.Now().UTC().Add(s.refreshTTL),
	}

	if err = s.sessionRepository.Create(ctx, session, s.refreshTTL); err != nil {
		return model.TokenPair{}, err
	}

	accessToken, err := s.generateAccessToken(user.UUID, session.UUID, time.Now().Add(s.accessTTL))
	if err != nil {
		return model.TokenPair{}, err
	}

	// Refresh-токеном служит UUID сессии: он ничего не значит сам по себе,
	// пока живёт соответствующая запись в Redis.
	return model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: session.UUID,
	}, nil
}
