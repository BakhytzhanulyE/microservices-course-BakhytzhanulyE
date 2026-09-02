// Package service объявляет бизнес-логику каталога деталей.
package service

import (
	"context"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/model"
)

// PartService — операции над каталогом деталей.
type PartService interface {
	Get(ctx context.Context, uuid string) (model.Part, error)
	List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
