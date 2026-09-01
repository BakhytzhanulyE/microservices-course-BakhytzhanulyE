package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/api/order/v1"
	grpcClient "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/client/grpc"
	inventoryClient "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/client/grpc/inventory"
	paymentClient "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/client/grpc/payment"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/config"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/repository"
	orderRepository "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/repository/order"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/service"
	orderService "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/service/order"
	orderProducer "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/service/producer/order_producer"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/closer"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka"
	kafkaProducer "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka/producer"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
	inventoryV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/payment/v1"
)

// diContainer лениво собирает зависимости сервиса заказов.
type diContainer struct {
	orderV1API *orderV1API.API

	orderService         service.OrderService
	orderProducerService service.OrderProducerService

	orderRepository repository.OrderRepository

	inventoryClient grpcClient.InventoryClient
	paymentClient   grpcClient.PaymentClient

	pgPool *pgxpool.Pool

	syncProducer      sarama.SyncProducer
	orderPaidProducer kafka.Producer
}

// NewDiContainer создаёт пустой контейнер зависимостей.
func NewDiContainer() *diContainer {
	return &diContainer{}
}

// OrderV1API возвращает HTTP-слой заказов.
func (d *diContainer) OrderV1API(ctx context.Context) *orderV1API.API {
	if d.orderV1API == nil {
		d.orderV1API = orderV1API.NewAPI(d.OrderService(ctx))
	}

	return d.orderV1API
}

// OrderService возвращает бизнес-логику заказов.
func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
			d.OrderProducerService(),
		)
	}

	return d.orderService
}

// OrderProducerService возвращает продюсер доменных событий.
func (d *diContainer) OrderProducerService() service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = orderProducer.NewService(d.OrderPaidProducer())
	}

	return d.orderProducerService
}

// OrderRepository возвращает хранилище заказов.
func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(d.PgPool(ctx))
	}

	return d.orderRepository
}

// PgPool создаёт пул соединений с PostgreSQL.
func (d *diContainer) PgPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.DSN())
		if err != nil {
			panic(fmt.Sprintf("не удалось создать пул PostgreSQL: %v", err))
		}

		if err = pool.Ping(ctx); err != nil {
			panic(fmt.Sprintf("PostgreSQL не отвечает на ping: %v", err))
		}

		closer.AddNamed("PostgreSQL pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}

	return d.pgPool
}

// InventoryClient создаёт gRPC-клиент каталога деталей.
func (d *diContainer) InventoryClient(_ context.Context) grpcClient.InventoryClient {
	if d.inventoryClient == nil {
		conn := d.dialGRPC("inventory", config.AppConfig().Clients.InventoryAddress())
		d.inventoryClient = inventoryClient.NewClient(inventoryV1.NewInventoryServiceClient(conn))
	}

	return d.inventoryClient
}

// PaymentClient создаёт gRPC-клиент платёжного сервиса.
func (d *diContainer) PaymentClient(_ context.Context) grpcClient.PaymentClient {
	if d.paymentClient == nil {
		conn := d.dialGRPC("payment", config.AppConfig().Clients.PaymentAddress())
		d.paymentClient = paymentClient.NewClient(paymentV1.NewPaymentServiceClient(conn))
	}

	return d.paymentClient
}

// dialGRPC создаёт соединение. grpc.NewClient не ходит в сеть сразу —
// соединение поднимется лениво при первом вызове, поэтому порядок старта сервисов не важен.
func (d *diContainer) dialGRPC(name, address string) *grpc.ClientConn {
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		panic(fmt.Sprintf("не удалось создать gRPC-клиент %s: %v", name, err))
	}

	closer.AddNamed(fmt.Sprintf("gRPC client %s", name), func(_ context.Context) error {
		return conn.Close()
	})

	return conn
}

// SyncProducer создаёт синхронный продюсер Kafka.
func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidProducer.Config(),
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

// OrderPaidProducer возвращает продюсер в топик order.paid.
func (d *diContainer) OrderPaidProducer() kafka.Producer {
	if d.orderPaidProducer == nil {
		d.orderPaidProducer = kafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidProducer.Topic(),
			logger.Instance(),
		)
	}

	return d.orderPaidProducer
}
