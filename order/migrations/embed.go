// Package migrations хранит SQL-миграции сервиса заказов и отдаёт их как embed.FS,
// чтобы бинарник накатывал схему сам и не тянул за собой каталог с файлами.
package migrations

import "embed"

// FS — встроенные в бинарник миграции.
//
//go:embed *.sql
var FS embed.FS
