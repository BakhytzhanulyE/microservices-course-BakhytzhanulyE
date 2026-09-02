// Package part реализует бизнес-логику каталога деталей.
package part

import (
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/repository"
	def "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/service"
)

var _ def.PartService = (*service)(nil)

type service struct {
	partRepository repository.PartRepository
}

// NewService создаёт сервис поверх хранилища деталей.
func NewService(partRepository repository.PartRepository) *service {
	return &service{
		partRepository: partRepository,
	}
}
