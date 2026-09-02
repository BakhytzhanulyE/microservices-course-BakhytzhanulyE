package notify

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/notification/internal/model"
)

type fakeSender struct {
	messages []string
	err      error
}

func (f *fakeSender) SendMessage(_ context.Context, text string) error {
	f.messages = append(f.messages, text)
	return f.err
}

func TestOrderPaid(t *testing.T) {
	t.Parallel()

	sender := &fakeSender{}
	s := NewService(sender)

	err := s.OrderPaid(context.Background(), model.OrderPaidEvent{
		OrderUUID:       "order-1",
		UserUUID:        "user-1",
		PaymentAmount:   150.5,
		TransactionUUID: "tx-1",
		PaymentMethod:   "CARD",
		PaidAt:          time.Now(),
	})

	require.NoError(t, err)
	require.Len(t, sender.messages, 1)

	text := sender.messages[0]
	assert.Contains(t, text, "order-1")
	assert.Contains(t, text, "150.50")
	assert.Contains(t, text, "CARD")
	assert.Contains(t, text, "tx-1")
}

func TestShipAssembled(t *testing.T) {
	t.Parallel()

	sender := &fakeSender{}
	s := NewService(sender)

	err := s.ShipAssembled(context.Background(), model.ShipAssembledEvent{
		OrderUUID:    "order-2",
		BuildTimeSec: 10,
	})

	require.NoError(t, err)
	require.Len(t, sender.messages, 1)
	assert.Contains(t, sender.messages[0], "order-2")
	assert.Contains(t, sender.messages[0], "10")
}

func TestSenderErrorIsPropagated(t *testing.T) {
	t.Parallel()

	sender := &fakeSender{err: errors.New("telegram down")}
	s := NewService(sender)

	require.Error(t, s.OrderPaid(context.Background(), model.OrderPaidEvent{OrderUUID: "order-3"}))
}
