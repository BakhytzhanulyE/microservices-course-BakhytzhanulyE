package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/client/http/telegram"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/config"
	kafkaConverter "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/converter/kafka"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/converter/kafka/decoder"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/service"
	orderConsumer "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/service/consumer/order_consumer"
	shipConsumer "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/service/consumer/ship_consumer"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/service/notify"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/closer"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka"
	wrappedConsumer "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka/consumer"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
	kafkaMiddleware "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/middleware/kafka"
)

// diContainer лениво собирает зависимости сервиса уведомлений.
type diContainer struct {
	orderConsumerService service.ConsumerService
	shipConsumerService  service.ConsumerService
	notifyService        service.NotifyService

	telegramClient *telegram.Client

	orderPaidConsumer     kafka.Consumer
	shipAssembledConsumer kafka.Consumer

	orderPaidDecoder     kafkaConverter.OrderPaidDecoder
	shipAssembledDecoder kafkaConverter.ShipAssembledDecoder
}

// NewDiContainer создаёт пустой контейнер зависимостей.
func NewDiContainer() *diContainer {
	return &diContainer{}
}

// OrderConsumerService возвращает консьюмер события OrderPaid.
func (d *diContainer) OrderConsumerService() service.ConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = orderConsumer.NewService(
			d.OrderPaidConsumer(),
			d.OrderPaidDecoder(),
			d.NotifyService(),
		)
	}

	return d.orderConsumerService
}

// ShipConsumerService возвращает консьюмер события ShipAssembled.
func (d *diContainer) ShipConsumerService() service.ConsumerService {
	if d.shipConsumerService == nil {
		d.shipConsumerService = shipConsumer.NewService(
			d.ShipAssembledConsumer(),
			d.ShipAssembledDecoder(),
			d.NotifyService(),
		)
	}

	return d.shipConsumerService
}

// NotifyService возвращает сервис уведомлений.
func (d *diContainer) NotifyService() service.NotifyService {
	if d.notifyService == nil {
		d.notifyService = notify.NewService(d.TelegramClient())
	}

	return d.notifyService
}

// TelegramClient возвращает клиента Telegram.
func (d *diContainer) TelegramClient() *telegram.Client {
	if d.telegramClient == nil {
		d.telegramClient = telegram.NewClient(
			config.AppConfig().Telegram.Token(),
			config.AppConfig().Telegram.ChatID(),
		)
	}

	return d.telegramClient
}

// OrderPaidConsumer возвращает консьюмер топика order.paid.
func (d *diContainer) OrderPaidConsumer() kafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = d.newConsumer("order.paid", config.AppConfig().OrderPaidConsumer)
	}

	return d.orderPaidConsumer
}

// ShipAssembledConsumer возвращает консьюмер топика ship.assembled.
func (d *diContainer) ShipAssembledConsumer() kafka.Consumer {
	if d.shipAssembledConsumer == nil {
		d.shipAssembledConsumer = d.newConsumer("ship.assembled", config.AppConfig().ShipAssembledConsumer)
	}

	return d.shipAssembledConsumer
}

// newConsumer поднимает отдельную consumer group на топик.
// Группы разные: иначе два консьюмера в одном сервисе делили бы партиции
// между собой и каждый видел бы только часть событий.
func (d *diContainer) newConsumer(name string, cfg config.ConsumerConfig) kafka.Consumer {
	group, err := sarama.NewConsumerGroup(config.AppConfig().Kafka.Brokers(), cfg.GroupID(), cfg.Config())
	if err != nil {
		panic(fmt.Sprintf("не удалось создать consumer group %s: %v", name, err))
	}

	closer.AddNamed(fmt.Sprintf("Kafka consumer group %s", name), func(_ context.Context) error {
		return group.Close()
	})

	return wrappedConsumer.NewConsumer(
		group,
		[]string{cfg.Topic()},
		logger.Instance(),
		kafkaMiddleware.Tracing(),
		kafkaMiddleware.Logging(logger.Instance()),
	)
}

// OrderPaidDecoder возвращает декодер события OrderPaid.
func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

// ShipAssembledDecoder возвращает декодер события ShipAssembled.
func (d *diContainer) ShipAssembledDecoder() kafkaConverter.ShipAssembledDecoder {
	if d.shipAssembledDecoder == nil {
		d.shipAssembledDecoder = decoder.NewShipAssembledDecoder()
	}

	return d.shipAssembledDecoder
}
