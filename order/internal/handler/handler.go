package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/storage"
)

// Handler обслуживает HTTP-запросы к заказам и хранит зависимости, которые для этого нужны.
type Handler struct {
	storage *storage.Storage
}

// NewHandler создаёт Handler поверх переданного хранилища.
func NewHandler(s *storage.Storage) *Handler {
	return &Handler{storage: s}
}

// Health отвечает 200 OK — простейшая проверка живости сервиса.
func (*Handler) Health(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte("OK")); err != nil {
		slog.Error("health: не удалось записать ответ", "error", err)
	}
}

// GetOrder отдаёт заказ по order_uuid из URL; 404, если такого нет.
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderUUID := chi.URLParam(r, "order_uuid")
	order, ok := h.storage.Get(orderUUID)
	if !ok {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		slog.Error("GetOrder: не удалось записать ответ", "error", err)
	}
}

func (h *Handler) PayOrder(w http.ResponseWriter, r *http.Request) {
	orderUUID := chi.URLParam(r, "order_uuid")

	var req model.PayOrderRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.storage.Pay(orderUUID)
	switch {
	case errors.Is(err, storage.ErrOrderNotFound):
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	case errors.Is(err, storage.ErrOrderNotPendingPayment):
		http.Error(w, "order is not awaiting payment", http.StatusConflict)
		return
	case err != nil:
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp := model.PayOrderResponse{
		TransactionUUID: uuid.NewString(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("PayOrder: не удалось записать ответ", "error", err)
		return
	}
}

// CreateOrder создаёт заказ из JSON-тела запроса и возвращает его с кодом 201.
// На некорректное тело отвечает 400.
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req model.CreateOrderRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	order := model.Order{
		UUID:       uuid.NewString(),
		UserUUID:   req.UserUUID,
		TotalPrice: req.TotalPrice,
		Status:     model.StatusPendingPayment,
	}

	h.storage.Create(order)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(order); err != nil {
		slog.Error("CreateOrder: не удалось записать ответ", "error", err)
	}
}
