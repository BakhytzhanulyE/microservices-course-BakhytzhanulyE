// Package part — реализация хранилища деталей поверх MongoDB.
package part

import (
	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/repository"
)

// Проверка на этапе компиляции, что repository реализует контракт.
var _ def.PartRepository = (*repository)(nil)

const collectionName = "parts"

type repository struct {
	collection *mongo.Collection
}

// NewRepository создаёт хранилище поверх коллекции parts.
func NewRepository(db *mongo.Database) *repository {
	return &repository{
		collection: db.Collection(collectionName),
	}
}
