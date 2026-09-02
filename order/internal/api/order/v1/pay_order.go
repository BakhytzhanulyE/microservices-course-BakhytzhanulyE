package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	paymentClient "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/client/grpc/payment"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
)

// PayOrder оплачивает заказ и возвращает UUID транзакции.
func (a *API) PayOrder(w http.ResponseWriter, r *http.Request) {
	orderUUID := chi.URLParam(r, "order_uuid")

	var req PayOrderRequest

	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, r, http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})

		return
	}

	// Способ оплаты проверяем здесь, чтобы не ходить в платёжный сервис зря.
	if !paymentClient.IsSupportedMethod(req.PaymentMethod) {
		writeError(w, r, model.ErrUnsupportedPaymentMethod)
		return
	}

	transactionUUID, err := a.orderService.Pay(r.Context(), model.PayOrderParams{
		OrderUUID:     orderUUID,
		PaymentMethod: req.PaymentMethod,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, r, http.StatusOK, PayOrderResponse{TransactionUUID: transactionUUID})
}
