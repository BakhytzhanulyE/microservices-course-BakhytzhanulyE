package model

type Order struct {
	UUID       string  `json:"order_uuid"`
	UserUUID   string  `json:"user_uuid"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
}

type CreateOrderRequest struct {
	UserUUID   string  `json:"user_uuid"`
	TotalPrice float64 `json:"total_price"`
}
