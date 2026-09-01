package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// CancelOrder отменяет неоплаченный заказ.
func (a *API) CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderUUID := chi.URLParam(r, "order_uuid")

	if err := a.orderService.Cancel(r.Context(), orderUUID); err != nil {
		writeError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
