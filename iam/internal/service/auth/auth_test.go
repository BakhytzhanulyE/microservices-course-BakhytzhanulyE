package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/model"
)

const (
	testSecret     = "test-secret-key"
	testAccessTTL  = 15 * time.Minute
	testRefreshTTL = 24 * time.Hour
)

type fakeUserRepo struct {
	users     map[string]model.User
	createErr error
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{users: make(map[string]model.User)}
}

func (f *fakeUserRepo) Create(_ context.Context, user model.User) error {
	if f.createErr != nil {
		return f.createErr
	}

	for _, existing := range f.users {
		if existing.Login == user.Login {
			return model.ErrUserAlreadyExists
		}
	}

	f.users[user.UUID] = user

	return nil
}

func (f *fakeUserRepo) GetByUUID(_ context.Context, uuid string) (model.User, error) {
	user, ok := f.users[uuid]
	if !ok {
		return model.User{}, model.ErrUserNotFound
	}

	return user, nil
}

func (f *fakeUserRepo) GetByLogin(_ context.Context, login string) (model.User, error) {
	for _, user := range f.users {
		if user.Login == login {
			return user, nil
		}
	}

	return model.User{}, model.ErrUserNotFound
}

type fakeSessionRepo struct {
	sessions map[string]model.Session
}

func newFakeSessionRepo() *fakeSessionRepo {
	return &fakeSessionRepo{sessions: make(map[string]model.Session)}
}

func (f *fakeSessionRepo) Create(_ context.Context, session model.Session, _ time.Duration) error {
	f.sessions[session.UUID] = session
	return nil
}

func (f *fakeSessionRepo) Get(_ context.Context, uuid string) (model.Session, error) {
	session, ok := f.sessions[uuid]
	if !ok {
		return model.Session{}, model.ErrSessionNotFound
	}

	return session, nil
}

func (f *fakeSessionRepo) Delete(_ context.Context, uuid string) error {
	delete(f.sessions, uuid)
	return nil
}

func newTestService(accessTTL time.Duration) (*service, *fakeUserRepo, *fakeSessionRepo) {
	userRepo := newFakeUserRepo()
	sessionRepo := newFakeSessionRepo()

	return NewService(userRepo, sessionRepo, testSecret, accessTTL, testRefreshTTL), userRepo, sessionRepo
}

func TestRegister(t *testing.T) {
	t.Parallel()

	t.Run("пароль сохраняется только хешем", func(t *testing.T) {
		t.Parallel()

		s, userRepo, _ := newTestService(testAccessTTL)

		userUUID, err := s.Register(context.Background(), model.RegisterParams{
			Login:    "ernat",
			Email:    "ernat@example.com",
			Password: "super-secret",
		})

		require.NoError(t, err)
		require.NotEmpty(t, userUUID)

		stored := userRepo.users[userUUID]
		assert.NotEqual(t, "super-secret", stored.PasswordHash)
		assert.NotEmpty(t, stored.PasswordHash)
	})

	t.Run("короткий пароль отклоняется", func(t *testing.T) {
		t.Parallel()

		s, _, _ := newTestService(testAccessTTL)

		_, err := s.Register(context.Background(), model.RegisterParams{
			Login:    "ernat",
			Email:    "ernat@example.com",
			Password: "123",
		})

		require.ErrorIs(t, err, model.ErrWeakPassword)
	})

	t.Run("занятый логин отклоняется", func(t *testing.T) {
		t.Parallel()

		s, _, _ := newTestService(testAccessTTL)
		params := model.RegisterParams{Login: "ernat", Email: "a@example.com", Password: "super-secret"}

		_, err := s.Register(context.Background(), params)
		require.NoError(t, err)

		_, err = s.Register(context.Background(), params)
		require.ErrorIs(t, err, model.ErrUserAlreadyExists)
	})
}

func TestLoginAndWhoami(t *testing.T) {
	t.Parallel()

	register := func(t *testing.T, s *service) {
		t.Helper()

		_, err := s.Register(context.Background(), model.RegisterParams{
			Login:    "ernat",
			Email:    "ernat@example.com",
			Password: "super-secret",
		})
		require.NoError(t, err)
	}

	t.Run("вход выдаёт токены, Whoami их разбирает", func(t *testing.T) {
		t.Parallel()

		s, _, _ := newTestService(testAccessTTL)
		register(t, s)

		tokens, err := s.Login(context.Background(), model.LoginParams{Login: "ernat", Password: "super-secret"})
		require.NoError(t, err)
		require.NotEmpty(t, tokens.AccessToken)
		require.NotEmpty(t, tokens.RefreshToken)

		session, user, err := s.Whoami(context.Background(), tokens.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, "ernat", user.Login)
		assert.Equal(t, tokens.RefreshToken, session.UUID)
	})

	t.Run("неверный пароль неотличим от несуществующего пользователя", func(t *testing.T) {
		t.Parallel()

		s, _, _ := newTestService(testAccessTTL)
		register(t, s)

		_, wrongPassErr := s.Login(context.Background(), model.LoginParams{Login: "ernat", Password: "не-тот"})
		_, noUserErr := s.Login(context.Background(), model.LoginParams{Login: "нет-такого", Password: "super-secret"})

		require.ErrorIs(t, wrongPassErr, model.ErrInvalidCredentials)
		require.ErrorIs(t, noUserErr, model.ErrInvalidCredentials)
	})

	t.Run("удалённая сессия делает живой токен бесполезным", func(t *testing.T) {
		t.Parallel()

		s, _, sessionRepo := newTestService(testAccessTTL)
		register(t, s)

		tokens, err := s.Login(context.Background(), model.LoginParams{Login: "ernat", Password: "super-secret"})
		require.NoError(t, err)

		require.NoError(t, sessionRepo.Delete(context.Background(), tokens.RefreshToken))

		_, _, err = s.Whoami(context.Background(), tokens.AccessToken)
		require.ErrorIs(t, err, model.ErrInvalidToken)
	})

	t.Run("протухший токен не проходит", func(t *testing.T) {
		t.Parallel()

		s, _, _ := newTestService(-time.Minute)
		register(t, s)

		tokens, err := s.Login(context.Background(), model.LoginParams{Login: "ernat", Password: "super-secret"})
		require.NoError(t, err)

		_, _, err = s.Whoami(context.Background(), tokens.AccessToken)
		require.ErrorIs(t, err, model.ErrInvalidToken)
	})

	t.Run("токен, подписанный чужим ключом, не проходит", func(t *testing.T) {
		t.Parallel()

		s, _, _ := newTestService(testAccessTTL)
		register(t, s)

		tokens, err := s.Login(context.Background(), model.LoginParams{Login: "ernat", Password: "super-secret"})
		require.NoError(t, err)

		other := NewService(newFakeUserRepo(), newFakeSessionRepo(), "другой-ключ", testAccessTTL, testRefreshTTL)

		_, _, err = other.Whoami(context.Background(), tokens.AccessToken)
		require.ErrorIs(t, err, model.ErrInvalidToken)
	})
}
