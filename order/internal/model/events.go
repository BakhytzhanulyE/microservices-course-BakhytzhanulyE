package model

import "time"

// OrderPaidEvent — событие «заказ оплачен», которое уезжает в Kafka.
type OrderPaidEvent struct {
	EventUUID       string
	OrderUUID       string
	UserUUID        string
	PaymentAmount   float64
	TransactionUUID string
	PaymentMethod   string
	PaidAt          time.Time
}
