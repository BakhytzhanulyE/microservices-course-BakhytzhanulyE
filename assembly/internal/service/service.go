// Package service объявляет бизнес-логику сборочного цеха.
package service

import (
	"context"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/assembly/internal/model"
)

// ConsumerService читает события из Kafka до отмены контекста.
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

// ShipProducerService публикует событие о собранном корабле.
type ShipProducerService interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}
