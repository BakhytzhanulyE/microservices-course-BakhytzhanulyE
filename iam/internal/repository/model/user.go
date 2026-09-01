// Package model описывает, как пользователь лежит в PostgreSQL.
package model

import "time"

// User — строка таблицы users.
type User struct {
	UUID         string    `db:"uuid"`
	Login        string    `db:"login"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}
