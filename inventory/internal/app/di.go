package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	inventoryV1API "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/api/inventory/v1"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/config"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/repository"
	partRepository "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/repository/part"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/service"
	partService "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/service/part"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/closer"
	inventoryV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/inventory/v1"
)

// diContainer собирает зависимости лениво: объект создаётся при первом обращении
// и переиспользуется дальше.
type diContainer struct {
	inventoryV1API inventoryV1.InventoryServiceServer

	partService    service.PartService
	partRepository repository.PartRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

// NewDiContainer создаёт пустой контейнер зависимостей.
func NewDiContainer() *diContainer {
	return &diContainer{}
}

// InventoryV1API возвращает gRPC-обработчик каталога.
func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryV1API.NewAPI(d.PartService(ctx))
	}

	return d.inventoryV1API
}

// PartService возвращает бизнес-логику каталога.
func (d *diContainer) PartService(ctx context.Context) service.PartService {
	if d.partService == nil {
		d.partService = partService.NewService(d.PartRepository(ctx))
	}

	return d.partService
}

// PartRepository возвращает хранилище деталей.
func (d *diContainer) PartRepository(ctx context.Context) repository.PartRepository {
	if d.partRepository == nil {
		d.partRepository = partRepository.NewRepository(d.MongoDBHandle(ctx))
	}

	return d.partRepository
}

// MongoDBClient подключается к MongoDB и регистрирует закрытие соединения.
func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("не удалось подключиться к MongoDB: %v", err))
		}

		if err = client.Ping(ctx, readpref.Primary()); err != nil {
			panic(fmt.Sprintf("MongoDB не отвечает на ping: %v", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

// MongoDBHandle возвращает базу данных сервиса.
func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}
