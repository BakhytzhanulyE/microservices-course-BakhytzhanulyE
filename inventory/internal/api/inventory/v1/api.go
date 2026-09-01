// Package v1 — gRPC-слой каталога деталей.
package v1

import (
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/service"
	inventoryV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer

	partService service.PartService
}

// NewAPI создаёт gRPC-обработчик каталога деталей.
func NewAPI(partService service.PartService) *api {
	return &api{
		partService: partService,
	}
}
