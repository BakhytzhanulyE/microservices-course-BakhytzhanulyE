// Package model описывает, как заказ лежит в PostgreSQL.
package model

import "time"

// Order — строка таблицы orders.
type Order struct {
	UUID            string     `db:"uuid"`
	UserUUID        string     `db:"user_uuid"`
	PartUUIDs       []string   `db:"part_uuids"`
	TotalPrice      float64    `db:"total_price"`
	TransactionUUID *string    `db:"transaction_uuid"`
	PaymentMethod   *string    `db:"payment_method"`
	Status          string     `db:"status"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
}
