// Package migrations хранит SQL-миграции IAM и отдаёт их как embed.FS.
package migrations

import "embed"

// FS — встроенные в бинарник миграции.
//
//go:embed *.sql
var FS embed.FS
