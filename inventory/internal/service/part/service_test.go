package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/model"
)

type fakeRepo struct {
	parts      []model.Part
	err        error
	lastFilter model.PartsFilter
}

func (f *fakeRepo) Get(_ context.Context, uuid string) (model.Part, error) {
	if f.err != nil {
		return model.Part{}, f.err
	}

	for _, part := range f.parts {
		if part.UUID == uuid {
			return part, nil
		}
	}

	return model.Part{}, model.ErrPartNotFound
}

func (f *fakeRepo) List(_ context.Context, filter model.PartsFilter) ([]model.Part, error) {
	f.lastFilter = filter
	return f.parts, f.err
}

func TestGet(t *testing.T) {
	t.Parallel()

	repo := &fakeRepo{parts: []model.Part{{UUID: "p1", Name: "Двигатель"}}}
	s := NewService(repo)

	t.Run("деталь находится по UUID", func(t *testing.T) {
		t.Parallel()

		part, err := s.Get(context.Background(), "p1")

		require.NoError(t, err)
		assert.Equal(t, "Двигатель", part.Name)
	})

	t.Run("неизвестный UUID отдаёт ErrPartNotFound", func(t *testing.T) {
		t.Parallel()

		_, err := s.Get(context.Background(), "нет-такой")

		require.ErrorIs(t, err, model.ErrPartNotFound)
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	repo := &fakeRepo{parts: []model.Part{{UUID: "p1"}, {UUID: "p2"}}}
	s := NewService(repo)

	filter := model.PartsFilter{Categories: []model.Category{model.CategoryEngine}}

	parts, err := s.List(context.Background(), filter)

	require.NoError(t, err)
	assert.Len(t, parts, 2)
	assert.Equal(t, filter, repo.lastFilter, "фильтр должен доезжать до хранилища без изменений")
}
