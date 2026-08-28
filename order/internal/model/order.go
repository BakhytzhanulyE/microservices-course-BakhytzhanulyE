package model

type Order struct {
	UUID       string  `json:"order_uuid"`
	UserUUID   string  `json:"user_uuid"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
}

const (
	StatusPendingPayment = "PENDING_PAYMENT"
	StatusPaid           = "PAID"
	StatusCompleted      = "COMPLETED"
	StatusShipped        = "SHIPPED"
)

type CreateOrderRequest struct {
	UserUUID   string  `json:"user_uuid"`
	TotalPrice float64 `json:"total_price"`
}

type PayOrderRequest struct {
	PaymentMethod string `json:"payment_method"`
}

type PayOrderResponse struct {
	TransactionUUID string `json:"transaction_uuid"`
}
