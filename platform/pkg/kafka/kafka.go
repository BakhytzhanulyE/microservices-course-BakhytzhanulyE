// Package kafka описывает интерфейсы продюсера и консьюмера, за которыми прячется sarama.
package kafka

import "context"

// MessageHandler — обработчик одного сообщения из топика.
type MessageHandler func(ctx context.Context, msg Message) error

// Consumer читает сообщения из топиков и отдаёт их обработчику.
type Consumer interface {
	Consume(ctx context.Context, handler MessageHandler) error
}

// Producer отправляет сообщения в топик.
type Producer interface {
	Send(ctx context.Context, key, value []byte) error
}
