// Package kafka объявляет декодеры входящих событий.
package kafka

import "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/model"

// OrderPaidDecoder разбирает сообщение в событие OrderPaid.
type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}

// ShipAssembledDecoder разбирает сообщение в событие ShipAssembled.
type ShipAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembledEvent, error)
}
