package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
)

// writeJSON отдаёт тело ответа в JSON с нужным статусом.
func writeJSON(w http.ResponseWriter, r *http.Request, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if payload == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logger.Error(r.Context(), "Не удалось записать тело ответа", zap.Error(err))
	}
}

// writeError переводит ошибку в HTTP-статус. Всё, что не разобрано явно,
// уходит как 500 — наружу детали внутренних сбоев не отдаём.
func writeError(w http.ResponseWriter, r *http.Request, err error) {
	status, message := http.StatusInternalServerError, "internal server error"

	switch {
	case errors.Is(err, model.ErrOrderNotFound):
		status, message = http.StatusNotFound, "order not found"
	case errors.Is(err, model.ErrPartsNotFound):
		status, message = http.StatusBadRequest, "some parts not found in inventory"
	case errors.Is(err, model.ErrEmptyPartList):
		status, message = http.StatusBadRequest, "part_uuids must not be empty"
	case errors.Is(err, model.ErrUnsupportedPaymentMethod):
		status, message = http.StatusBadRequest, "unsupported payment method"
	case errors.Is(err, model.ErrOrderNotPendingPayment):
		status, message = http.StatusConflict, "order is not awaiting payment"
	case errors.Is(err, model.ErrOrderAlreadyPaid):
		status, message = http.StatusConflict, "paid order cannot be cancelled"
	case errors.Is(err, model.ErrInventoryUnavailable):
		status, message = http.StatusServiceUnavailable, "inventory service unavailable"
	case errors.Is(err, model.ErrPaymentFailed):
		status, message = http.StatusBadGateway, "payment service failed"
	}

	if status >= http.StatusInternalServerError {
		logger.Error(r.Context(), "Запрос завершился ошибкой", zap.Error(err))
	}

	writeJSON(w, r, status, ErrorResponse{Code: status, Message: message})
}

// decodeJSON разбирает тело запроса и запрещает неизвестные поля,
// чтобы опечатка в имени поля не проходила молча.
func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(dst)
}
