// Package model описывает события, на которые реагирует сервис уведомлений.
package model

import "time"

// OrderPaidEvent — событие «заказ оплачен».
type OrderPaidEvent struct {
	EventUUID       string
	OrderUUID       string
	UserUUID        string
	PaymentAmount   float64
	TransactionUUID string
	PaymentMethod   string
	PaidAt          time.Time
}

// ShipAssembledEvent — событие «корабль собран».
type ShipAssembledEvent struct {
	EventUUID    string
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
	AssembledAt  time.Time
}
