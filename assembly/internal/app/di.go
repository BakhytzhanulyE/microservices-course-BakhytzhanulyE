package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/assembly/internal/config"
	kafkaConverter "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/assembly/internal/converter/kafka"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/assembly/internal/converter/kafka/decoder"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/assembly/internal/service"
	orderConsumer "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/assembly/internal/service/consumer/order_consumer"
	shipProducer "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/assembly/internal/service/producer/ship_producer"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/closer"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka"
	wrappedConsumer "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka/consumer"
	wrappedProducer "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka/producer"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
	kafkaMiddleware "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/middleware/kafka"
)

// diContainer лениво собирает зависимости сборочного цеха.
type diContainer struct {
	consumerService     service.ConsumerService
	shipProducerService service.ShipProducerService

	consumerGroup     sarama.ConsumerGroup
	orderPaidConsumer kafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder

	syncProducer          sarama.SyncProducer
	shipAssembledProducer kafka.Producer
}

// NewDiContainer создаёт пустой контейнер зависимостей.
func NewDiContainer() *diContainer {
	return &diContainer{}
}

// ConsumerService возвращает консьюмер оплаченных заказов.
func (d *diContainer) ConsumerService() service.ConsumerService {
	if d.consumerService == nil {
		d.consumerService = orderConsumer.NewService(
			d.OrderPaidConsumer(),
			d.OrderPaidDecoder(),
			d.ShipProducerService(),
			config.AppConfig().OrderPaidConsumer.BuildDuration(),
		)
	}

	return d.consumerService
}

// ShipProducerService возвращает продюсер событий сборки.
func (d *diContainer) ShipProducerService() service.ShipProducerService {
	if d.shipProducerService == nil {
		d.shipProducerService = shipProducer.NewService(d.ShipAssembledProducer())
	}

	return d.shipProducerService
}

// ConsumerGroup создаёт consumer group Kafka.
func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		group, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("не удалось создать consumer group: %v", err))
		}

		closer.AddNamed("Kafka consumer group", func(_ context.Context) error {
			return group.Close()
		})

		d.consumerGroup = group
	}

	return d.consumerGroup
}

// OrderPaidConsumer возвращает консьюмер топика order.paid.
func (d *diContainer) OrderPaidConsumer() kafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{config.AppConfig().OrderPaidConsumer.Topic()},
			logger.Instance(),
			kafkaMiddleware.Logging(logger.Instance()),
		)
	}

	return d.orderPaidConsumer
}

// OrderPaidDecoder возвращает декодер события OrderPaid.
func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

// SyncProducer создаёт синхронный продюсер Kafka.
func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().ShipAssembledProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("не удалось создать Kafka-продюсер: %v", err))
		}

		closer.AddNamed("Kafka sync producer", func(_ context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

// ShipAssembledProducer возвращает продюсер в топик ship.assembled.
func (d *diContainer) ShipAssembledProducer() kafka.Producer {
	if d.shipAssembledProducer == nil {
		d.shipAssembledProducer = wrappedProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().ShipAssembledProducer.Topic(),
			logger.Instance(),
		)
	}

	return d.shipAssembledProducer
}
