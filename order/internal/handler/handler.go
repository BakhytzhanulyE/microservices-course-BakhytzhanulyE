package handler

import (
	"encoding/json"
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
	order := &model.Order{
		UUID:       uuid.NewString(),
		UserUUID:   req.UserUUID,
		TotalPrice: req.TotalPrice,
		Status:     "PENDING_PAYMENT",
	}

	h.storage.Create(order)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(order); err != nil {
		slog.Error("CreateOrder: не удалось записать ответ", "error", err)
	}
}
