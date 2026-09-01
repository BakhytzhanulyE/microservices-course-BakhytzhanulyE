package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/model"
)

// claims — полезная нагрузка access-токена.
// sid хранит UUID сессии: по нему проверяем, что сессию не отозвали.
type claims struct {
	jwt.RegisteredClaims

	SessionUUID string `json:"sid"`
}

// generateAccessToken подписывает access-токен для пары «пользователь + сессия».
func (s *service) generateAccessToken(userUUID, sessionUUID string, expiresAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userUUID,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		SessionUUID: sessionUUID,
	})

	return token.SignedString([]byte(s.secretKey))
}

// parseAccessToken проверяет подпись и срок жизни токена.
func (s *service) parseAccessToken(tokenString string) (claims, error) {
	var parsed claims

	// WithValidMethods обязателен: без него токен с alg=none пройдёт проверку.
	_, err := jwt.ParseWithClaims(tokenString, &parsed, func(_ *jwt.Token) (any, error) {
		return []byte(s.secretKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return claims{}, model.ErrInvalidToken
		}

		return claims{}, model.ErrInvalidToken
	}

	return parsed, nil
}
