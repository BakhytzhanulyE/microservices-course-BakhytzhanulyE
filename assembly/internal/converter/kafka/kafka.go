// Package kafka объявляет декодеры входящих событий.
package kafka

import "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/assembly/internal/model"

// OrderPaidDecoder разбирает тело сообщения в событие OrderPaid.
type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}
