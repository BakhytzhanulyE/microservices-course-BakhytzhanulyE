// Package repository объявляет контракт хранилища деталей.
package repository

import (
	"context"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/model"
)

// PartRepository — хранилище деталей.
type PartRepository interface {
	Get(ctx context.Context, uuid string) (model.Part, error)
	List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
