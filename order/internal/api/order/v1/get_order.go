package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetOrder отдаёт заказ по order_uuid из URL.
func (a *API) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderUUID := chi.URLParam(r, "order_uuid")

	order, err := a.orderService.Get(r.Context(), orderUUID)
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, OrderResponse{
		OrderUUID:       order.UUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          string(order.Status),
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	})
}
