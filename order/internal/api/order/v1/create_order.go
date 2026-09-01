package v1

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
)

// CreateOrder создаёт заказ и возвращает его UUID и сумму.
func (a *API) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest

	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, r, http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})

		return
	}

	if _, err := uuid.Parse(req.UserUUID); err != nil {
		writeJSON(w, r, http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "user_uuid must be a valid uuid",
		})

		return
	}

	order, err := a.orderService.Create(r.Context(), model.CreateOrderParams{
		UserUUID:  req.UserUUID,
		PartUUIDs: req.PartUUIDs,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusCreated, CreateOrderResponse{
		OrderUUID:  order.UUID,
		TotalPrice: order.TotalPrice,
	})
}
