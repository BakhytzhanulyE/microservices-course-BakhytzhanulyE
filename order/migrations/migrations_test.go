package migrations

import (
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
)

// Опечатка в имени файла или в маркерах +goose всплыла бы только при старте сервиса
// против живой базы. Тест ловит это без всякой базы.
func TestMigrationsAreCollectable(t *testing.T) {
	goose.SetBaseFS(FS)
	require.NoError(t, goose.SetDialect("postgres"))

	migrations, err := goose.CollectMigrations(".", 0, goose.MaxVersion)
	require.NoError(t, err)
	require.NotEmpty(t, migrations, "goose не нашёл ни одной миграции")
}
