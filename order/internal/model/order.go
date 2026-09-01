// Package model описывает доменные сущности заказа.
package model

import "time"

// OrderStatus — состояние заказа.
type OrderStatus string

const (
	// OrderStatusPendingPayment — заказ создан и ждёт оплаты.
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	// OrderStatusPaid — заказ оплачен.
	OrderStatusPaid OrderStatus = "PAID"
	// OrderStatusCancelled — заказ отменён.
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// Order — заказ на детали космического корабля.
type Order struct {
	UUID            string
	UserUUID        string
	PartUUIDs       []string
	TotalPrice      float64
	TransactionUUID *string
	PaymentMethod   *string
	Status          OrderStatus
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}

// CreateOrderParams — что нужно, чтобы создать заказ.
type CreateOrderParams struct {
	UserUUID  string
	PartUUIDs []string
}

// PayOrderParams — что нужно, чтобы оплатить заказ.
type PayOrderParams struct {
	OrderUUID     string
	PaymentMethod string
}

// UpdateOrderParams — поля, которые можно изменить у заказа.
// nil означает «не трогать это поле».
type UpdateOrderParams struct {
	Status          *OrderStatus
	TransactionUUID *string
	PaymentMethod   *string
}

// Part — деталь, как её видит сервис заказов. Полная карточка живёт в inventory.
type Part struct {
	UUID  string
	Name  string
	Price float64
}
