// Package migrator накатывает SQL-миграции на PostgreSQL через goose.
package migrator

import (
	"context"
	"database/sql"
	"io/fs"

	// Драйвер pgx в режиме database/sql — нужен goose, который работает через *sql.DB.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
)

// Up накатывает миграции из migrationsFS на базу по строке подключения dsn.
// dir — путь к миграциям внутри FS (обычно "migrations").
func Up(ctx context.Context, dsn string, migrationsFS fs.FS, dir string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			logger.Error(ctx, "Не удалось закрыть соединение для миграций", zap.Error(closeErr))
		}
	}()

	if err = db.PingContext(ctx); err != nil {
		return err
	}

	goose.SetBaseFS(migrationsFS)

	if err = goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.UpContext(ctx, db, dir)
}
