package auth

import (
	"context"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/model"
)

// Whoami возвращает сессию и пользователя по access-токену.
func (s *service) Whoami(ctx context.Context, accessToken string) (model.Session, model.User, error) {
	parsed, err := s.parseAccessToken(accessToken)
	if err != nil {
		return model.Session{}, model.User{}, err
	}

	// Подписи мало: сессию могли удалить при выходе, а токен ещё не протух.
	session, err := s.sessionRepository.Get(ctx, parsed.SessionUUID)
	if err != nil {
		return model.Session{}, model.User{}, model.ErrInvalidToken
	}

	user, err := s.userRepository.GetByUUID(ctx, session.UserUUID)
	if err != nil {
		return model.Session{}, model.User{}, err
	}

	return session, user, nil
}
