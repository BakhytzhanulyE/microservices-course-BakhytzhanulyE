package consumer

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka"
)

// Middleware оборачивает обработчик сообщения — логированием, метриками и т.п.
type Middleware func(next kafka.MessageHandler) kafka.MessageHandler

type groupHandler struct {
	handler kafka.MessageHandler
	logger  Logger
}

// NewGroupHandler навешивает цепочку middleware на обработчик и возвращает
// sarama.ConsumerGroupHandler.
func NewGroupHandler(handler kafka.MessageHandler, logger Logger, middlewares ...Middleware) sarama.ConsumerGroupHandler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return &groupHandler{
		handler: handler,
		logger:  logger,
	}
}

func (*groupHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

func (*groupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (g *groupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				g.logger.Info(session.Context(), "Канал сообщений Kafka закрыт")
				return nil
			}

			msg := kafka.Message{
				Key:            message.Key,
				Value:          message.Value,
				Topic:          message.Topic,
				Partition:      message.Partition,
				Offset:         message.Offset,
				Timestamp:      message.Timestamp,
				BlockTimestamp: message.BlockTimestamp,
				Headers:        extractHeaders(message.Headers),
			}

			// Сообщение помечаем только после успешной обработки: иначе при падении
			// обработчика оффсет уедет вперёд и событие потеряется.
			if err := g.handler(session.Context(), msg); err != nil {
				g.logger.Error(session.Context(), "Ошибка обработчика Kafka", zap.Error(err))
				continue
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			g.logger.Info(session.Context(), "Контекст сессии Kafka завершён")
			return nil
		}
	}
}

func extractHeaders(headers []*sarama.RecordHeader) map[string][]byte {
	result := make(map[string][]byte, len(headers))
	for _, h := range headers {
		if h != nil && h.Key != nil {
			result[string(h.Key)] = h.Value
		}
	}

	return result
}
