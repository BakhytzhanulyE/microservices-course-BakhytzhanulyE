package app

import (
	"context"
	"fmt"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"

	iamV1API "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/api/iam/v1"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/config"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/repository"
	sessionRepository "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/repository/session"
	userRepository "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/repository/user"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/service"
	authService "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/iam/internal/service/auth"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/cache"
	redisClient "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/cache/redis"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/closer"
	iamV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/iam/v1"
)

// diContainer лениво собирает зависимости IAM.
type diContainer struct {
	iamV1API iamV1.AuthServiceServer

	authService service.AuthService

	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository

	pgPool      *pgxpool.Pool
	cacheClient cache.Client
}

// NewDiContainer создаёт пустой контейнер зависимостей.
func NewDiContainer() *diContainer {
	return &diContainer{}
}

// IamV1API возвращает gRPC-обработчик аутентификации.
func (d *diContainer) IamV1API(ctx context.Context) iamV1.AuthServiceServer {
	if d.iamV1API == nil {
		d.iamV1API = iamV1API.NewAPI(d.AuthService(ctx))
	}

	return d.iamV1API
}

// AuthService возвращает бизнес-логику аутентификации.
func (d *diContainer) AuthService(ctx context.Context) service.AuthService {
	if d.authService == nil {
		jwtCfg := config.AppConfig().JWT

		d.authService = authService.NewService(
			d.UserRepository(ctx),
			d.SessionRepository(ctx),
			jwtCfg.SecretKey(),
			jwtCfg.AccessTTL(),
			jwtCfg.RefreshTTL(),
		)
	}

	return d.authService
}

// UserRepository возвращает хранилище пользователей.
func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userRepository.NewRepository(d.PgPool(ctx))
	}

	return d.userRepository
}

// SessionRepository возвращает хранилище сессий.
func (d *diContainer) SessionRepository(ctx context.Context) repository.SessionRepository {
	if d.sessionRepository == nil {
		d.sessionRepository = sessionRepository.NewRepository(d.CacheClient(ctx))
	}

	return d.sessionRepository
}

// PgPool создаёт пул соединений с PostgreSQL.
func (d *diContainer) PgPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		// ParseConfig + NewWithConfig, а не pgxpool.New: пулу нужно подсунуть
		// трейсер, иначе SQL-запросы не попадают в трейс и на графике видно
		// только «ручка заняла 77 мс» без ответа, сколько из них ушло в базу.
		poolConfig, err := pgxpool.ParseConfig(config.AppConfig().Postgres.DSN())
		if err != nil {
			panic(fmt.Sprintf("не удалось разобрать DSN PostgreSQL: %v", err))
		}

		// WithIncludeQueryParameters намеренно НЕ включаем: значения параметров
		// уехали бы в Jaeger, а среди них пароли и персональные данные.
		poolConfig.ConnConfig.Tracer = otelpgx.NewTracer(
			otelpgx.WithTrimSQLInSpanName(),
		)

		pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
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

// CacheClient создаёт клиента Redis.
func (d *diContainer) CacheClient(ctx context.Context) cache.Client {
	if d.cacheClient == nil {
		redisCfg := config.AppConfig().Redis

		client, err := redisClient.NewClient(redisCfg.Address(), redisCfg.Password(), redisCfg.DB())
		if err != nil {
			panic(fmt.Sprintf("не удалось создать клиент Redis: %v", err))
		}

		if err := client.Ping(ctx); err != nil {
			panic(fmt.Sprintf("Redis не отвечает на ping: %v", err))
		}

		closer.AddNamed("Redis client", func(_ context.Context) error {
			return client.Close()
		})

		d.cacheClient = client
	}

	return d.cacheClient
}
