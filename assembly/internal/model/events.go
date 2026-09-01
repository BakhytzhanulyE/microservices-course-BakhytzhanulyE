// Package model описывает события, с которыми работает сборочный цех.
package model

import "time"

// OrderPaidEvent — событие «заказ оплачен», приходит из Kafka.
type OrderPaidEvent struct {
	EventUUID       string
	OrderUUID       string
	UserUUID        string
	PaymentAmount   float64
	TransactionUUID string
	PaymentMethod   string
	PaidAt          time.Time
}

// ShipAssembledEvent — событие «корабль собран», уезжает в Kafka.
type ShipAssembledEvent struct {
	EventUUID    string
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
	AssembledAt  time.Time
}
