package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/repository"
)

// Заглушки пишем руками: интерфейсы маленькие, а тесты так читаются без генератора моков.

type fakeRepo struct {
	orders    map[string]model.Order
	createErr error
	getErr    error
	updateErr error

	createCalls int
	updateCalls int
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{orders: make(map[string]model.Order)}
}

func (f *fakeRepo) Create(_ context.Context, order model.Order) error {
	f.createCalls++
	if f.createErr != nil {
		return f.createErr
	}

	f.orders[order.UUID] = order

	return nil
}

func (f *fakeRepo) Get(_ context.Context, uuid string) (model.Order, error) {
	if f.getErr != nil {
		return model.Order{}, f.getErr
	}

	order, ok := f.orders[uuid]
	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	return order, nil
}

func (f *fakeRepo) Update(_ context.Context, uuid string, params model.UpdateOrderParams) (model.Order, error) {
	f.updateCalls++
	if f.updateErr != nil {
		return model.Order{}, f.updateErr
	}

	order, ok := f.orders[uuid]
	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	if params.Status != nil {
		order.Status = *params.Status
	}

	if params.TransactionUUID != nil {
		order.TransactionUUID = params.TransactionUUID
	}

	if params.PaymentMethod != nil {
		order.PaymentMethod = params.PaymentMethod
	}

	f.orders[uuid] = order

	return order, nil
}

var _ repository.OrderRepository = (*fakeRepo)(nil)

type fakeInventory struct {
	parts []model.Part
	err   error

	lastUUIDs []string
}

func (f *fakeInventory) ListParts(_ context.Context, uuids []string) ([]model.Part, error) {
	f.lastUUIDs = uuids
	return f.parts, f.err
}

type fakePayment struct {
	transactionUUID string
	err             error

	calls int
}

func (f *fakePayment) PayOrder(_ context.Context, _, _, _ string) (string, error) {
	f.calls++
	return f.transactionUUID, f.err
}

type fakeProducer struct {
	events []model.OrderPaidEvent
	err    error
}

func (f *fakeProducer) ProduceOrderPaid(_ context.Context, event model.OrderPaidEvent) error {
	f.events = append(f.events, event)
	return f.err
}

func TestCreate(t *testing.T) {
	t.Parallel()

	engine := model.Part{UUID: "part-1", Name: "Двигатель", Price: 100}
	wing := model.Part{UUID: "part-2", Name: "Крыло", Price: 50.5}

	tests := []struct {
		name       string
		partUUIDs  []string
		inventory  *fakeInventory
		wantErr    error
		wantTotal  float64
		wantUnique int
	}{
		{
			name:       "две детали складываются в сумму заказа",
			partUUIDs:  []string{"part-1", "part-2"},
			inventory:  &fakeInventory{parts: []model.Part{engine, wing}},
			wantTotal:  150.5,
			wantUnique: 2,
		},
		{
			name:       "дубли деталей считаются один раз",
			partUUIDs:  []string{"part-1", "part-1"},
			inventory:  &fakeInventory{parts: []model.Part{engine}},
			wantTotal:  100,
			wantUnique: 1,
		},
		{
			name:      "пустой список деталей отвергается",
			partUUIDs: nil,
			inventory: &fakeInventory{},
			wantErr:   model.ErrEmptyPartList,
		},
		{
			name:      "неизвестная деталь ломает заказ",
			partUUIDs: []string{"part-1", "part-unknown"},
			inventory: &fakeInventory{parts: []model.Part{engine}},
			wantErr:   model.ErrPartsNotFound,
		},
		{
			name:      "недоступный каталог отдаёт понятную ошибку",
			partUUIDs: []string{"part-1"},
			inventory: &fakeInventory{err: errors.New("connection refused")},
			wantErr:   model.ErrInventoryUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := newFakeRepo()
			s := NewService(repo, tt.inventory, &fakePayment{}, &fakeProducer{})

			order, err := s.Create(context.Background(), model.CreateOrderParams{
				UserUUID:  "user-1",
				PartUUIDs: tt.partUUIDs,
			})

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Zero(t, repo.createCalls, "при ошибке заказ сохранять нельзя")

				return
			}

			require.NoError(t, err)
			assert.InDelta(t, tt.wantTotal, order.TotalPrice, 0.001)
			assert.Equal(t, model.OrderStatusPendingPayment, order.Status)
			assert.Len(t, order.PartUUIDs, tt.wantUnique)
			assert.Equal(t, 1, repo.createCalls)
		})
	}
}

func TestPay(t *testing.T) {
	t.Parallel()

	const orderUUID = "order-1"

	newOrder := func(status model.OrderStatus) model.Order {
		return model.Order{
			UUID:       orderUUID,
			UserUUID:   "user-1",
			TotalPrice: 200,
			Status:     status,
			CreatedAt:  time.Now(),
		}
	}

	t.Run("успешная оплата меняет статус и шлёт событие", func(t *testing.T) {
		t.Parallel()

		repo := newFakeRepo()
		repo.orders[orderUUID] = newOrder(model.OrderStatusPendingPayment)

		payment := &fakePayment{transactionUUID: "tx-1"}
		producer := &fakeProducer{}

		s := NewService(repo, &fakeInventory{}, payment, producer)

		transactionUUID, err := s.Pay(context.Background(), model.PayOrderParams{
			OrderUUID:     orderUUID,
			PaymentMethod: "CARD",
		})

		require.NoError(t, err)
		assert.Equal(t, "tx-1", transactionUUID)
		assert.Equal(t, model.OrderStatusPaid, repo.orders[orderUUID].Status)

		require.Len(t, producer.events, 1)
		assert.Equal(t, orderUUID, producer.events[0].OrderUUID)
		assert.InDelta(t, 200.0, producer.events[0].PaymentAmount, 0.001)
	})

	t.Run("повторная оплата отклоняется без похода в платёжный сервис", func(t *testing.T) {
		t.Parallel()

		repo := newFakeRepo()
		repo.orders[orderUUID] = newOrder(model.OrderStatusPaid)

		payment := &fakePayment{transactionUUID: "tx-2"}
		s := NewService(repo, &fakeInventory{}, payment, &fakeProducer{})

		_, err := s.Pay(context.Background(), model.PayOrderParams{OrderUUID: orderUUID, PaymentMethod: "CARD"})

		require.ErrorIs(t, err, model.ErrOrderNotPendingPayment)
		assert.Zero(t, payment.calls)
	})

	t.Run("несуществующий заказ отдаёт ErrOrderNotFound", func(t *testing.T) {
		t.Parallel()

		s := NewService(newFakeRepo(), &fakeInventory{}, &fakePayment{}, &fakeProducer{})

		_, err := s.Pay(context.Background(), model.PayOrderParams{OrderUUID: "нет-такого", PaymentMethod: "CARD"})

		require.ErrorIs(t, err, model.ErrOrderNotFound)
	})

	t.Run("сбой платёжного сервиса не меняет статус заказа", func(t *testing.T) {
		t.Parallel()

		repo := newFakeRepo()
		repo.orders[orderUUID] = newOrder(model.OrderStatusPendingPayment)

		payment := &fakePayment{err: errors.New("payment down")}
		s := NewService(repo, &fakeInventory{}, payment, &fakeProducer{})

		_, err := s.Pay(context.Background(), model.PayOrderParams{OrderUUID: orderUUID, PaymentMethod: "CARD"})

		require.ErrorIs(t, err, model.ErrPaymentFailed)
		assert.Equal(t, model.OrderStatusPendingPayment, repo.orders[orderUUID].Status)
		assert.Zero(t, repo.updateCalls)
	})

	t.Run("сбой Kafka не отменяет успешную оплату", func(t *testing.T) {
		t.Parallel()

		repo := newFakeRepo()
		repo.orders[orderUUID] = newOrder(model.OrderStatusPendingPayment)

		producer := &fakeProducer{err: errors.New("kafka down")}
		s := NewService(repo, &fakeInventory{}, &fakePayment{transactionUUID: "tx-3"}, producer)

		transactionUUID, err := s.Pay(context.Background(), model.PayOrderParams{OrderUUID: orderUUID, PaymentMethod: "CARD"})

		require.NoError(t, err)
		assert.Equal(t, "tx-3", transactionUUID)
		assert.Equal(t, model.OrderStatusPaid, repo.orders[orderUUID].Status)
	})
}

func TestCancel(t *testing.T) {
	t.Parallel()

	const orderUUID = "order-1"

	tests := []struct {
		name    string
		status  model.OrderStatus
		wantErr error
	}{
		{name: "неоплаченный заказ отменяется", status: model.OrderStatusPendingPayment},
		{name: "оплаченный заказ отменить нельзя", status: model.OrderStatusPaid, wantErr: model.ErrOrderAlreadyPaid},
		{name: "повторная отмена отклоняется", status: model.OrderStatusCancelled, wantErr: model.ErrOrderNotPendingPayment},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := newFakeRepo()
			repo.orders[orderUUID] = model.Order{UUID: orderUUID, Status: tt.status}

			s := NewService(repo, &fakeInventory{}, &fakePayment{}, &fakeProducer{})

			err := s.Cancel(context.Background(), orderUUID)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, model.OrderStatusCancelled, repo.orders[orderUUID].Status)
		})
	}
}
