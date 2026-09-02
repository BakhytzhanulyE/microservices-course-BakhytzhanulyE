// Package v1 — HTTP-слой сервиса заказов.
package v1

import "time"

// CreateOrderRequest — тело запроса на создание заказа.
type CreateOrderRequest struct {
	UserUUID  string   `json:"user_uuid"`
	PartUUIDs []string `json:"part_uuids"`
}

// CreateOrderResponse — ответ на создание заказа.
type CreateOrderResponse struct {
	OrderUUID  string  `json:"order_uuid"`
	TotalPrice float64 `json:"total_price"`
}

// PayOrderRequest — тело запроса на оплату.
type PayOrderRequest struct {
	PaymentMethod string `json:"payment_method"`
}

// PayOrderResponse — ответ на оплату.
type PayOrderResponse struct {
	TransactionUUID string `json:"transaction_uuid"`
}

// OrderResponse — представление заказа в HTTP API.
type OrderResponse struct {
	OrderUUID       string     `json:"order_uuid"`
	UserUUID        string     `json:"user_uuid"`
	PartUUIDs       []string   `json:"part_uuids"`
	TotalPrice      float64    `json:"total_price"`
	TransactionUUID *string    `json:"transaction_uuid,omitempty"`
	PaymentMethod   *string    `json:"payment_method,omitempty"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}

// ErrorResponse — единый формат ошибки API.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
