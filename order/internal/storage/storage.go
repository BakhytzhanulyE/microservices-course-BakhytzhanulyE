package storage

import "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"

type Storage struct {
	orders map[string]model.Order
}

func NewStorage() *Storage {
	return &Storage{
		orders: map[string]model.Order{
			"1": {UUID: "1", UserUUID: "user1", TotalPrice: 100.0, Status: "pending"},
			"2": {UUID: "2", UserUUID: "user2", TotalPrice: 200.0, Status: "completed"},
			"3": {UUID: "3", UserUUID: "user3", TotalPrice: 300.0, Status: "shipped"},
		},
	}
}

func (s *Storage) Get(orderUUID string) (model.Order, bool) {
	order, ok := s.orders[orderUUID]
	return order, ok
}

func (s *Storage) Create(order model.Order) {
	s.orders[order.UUID] = order
}

func (s *Storage) Update(order model.Order) bool {
	_, ok := s.orders[order.UUID]

	if ok {
		s.orders[order.UUID] = order
	}

	return ok
}
