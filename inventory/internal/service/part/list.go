package part

import (
	"context"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/model"
)

// List возвращает детали по фильтру.
func (s *service) List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	return s.partRepository.List(ctx, filter)
}
