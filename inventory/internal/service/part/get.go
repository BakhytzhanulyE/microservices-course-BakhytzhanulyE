package part

import (
	"context"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/model"
)

// Get возвращает деталь по UUID.
func (s *service) Get(ctx context.Context, uuid string) (model.Part, error) {
	return s.partRepository.Get(ctx, uuid)
}
