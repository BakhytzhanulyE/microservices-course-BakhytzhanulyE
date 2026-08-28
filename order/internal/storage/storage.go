package storage

import (
	"errors"
	"sync"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
)

var (
	ErrOrderNotFound          = errors.New("order not found")
	ErrOrderNotPendingPayment = errors.New("order is not pending payment")
)

type Storage struct {
	mu     sync.RWMutex
	orders map[string]model.Order
}

func NewStorage() *Storage {
	return &Storage{
		orders: map[string]model.Order{
			"1": {UUID: "1", UserUUID: "user1", TotalPrice: 100.0, Status: model.StatusPendingPayment},
			"2": {UUID: "2", UserUUID: "user2", TotalPrice: 200.0, Status: model.StatusCompleted},
			"3": {UUID: "3", UserUUID: "user3", TotalPrice: 300.0, Status: model.StatusShipped},
		},
	}
}

func (s *Storage) Get(orderUUID string) (model.Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[orderUUID]
	return order, ok
}

func (s *Storage) Create(order model.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.UUID] = order
}

func (s *Storage) Pay(orderUUID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	order, ok := s.orders[orderUUID]
	if !ok {
		return ErrOrderNotFound
	}
	if order.Status != model.StatusPendingPayment {
		return ErrOrderNotPendingPayment
	}

	order.Status = model.StatusPaid

	s.orders[orderUUID] = order

	return nil
}
