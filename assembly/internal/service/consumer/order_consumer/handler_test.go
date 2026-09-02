package order_consumer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/assembly/internal/converter/kafka/decoder"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/assembly/internal/model"
	platformKafka "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/kafka"
	eventsV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/events/v1"
)

type fakeConsumer struct{}

func (*fakeConsumer) Consume(context.Context, platformKafka.MessageHandler) error { return nil }

type fakeShipProducer struct {
	events []model.ShipAssembledEvent
	err    error
}

func (f *fakeShipProducer) ProduceShipAssembled(_ context.Context, event model.ShipAssembledEvent) error {
	f.events = append(f.events, event)
	return f.err
}

func orderPaidMessage(t *testing.T, orderUUID string) platformKafka.Message {
	t.Helper()

	payload, err := proto.Marshal(&eventsV1.OrderPaid{
		EventUuid:       "event-1",
		OrderUuid:       orderUUID,
		UserUuid:        "user-1",
		PaymentAmount:   150.5,
		TransactionUuid: "tx-1",
		PaymentMethod:   "CARD",
		PaidAt:          timestamppb.New(time.Now()),
	})
	require.NoError(t, err)

	return platformKafka.Message{Topic: "order.paid", Value: payload}
}

func TestOrderPaidHandler(t *testing.T) {
	t.Parallel()

	t.Run("после сборки публикуется ShipAssembled", func(t *testing.T) {
		t.Parallel()

		producer := &fakeShipProducer{}
		s := NewService(&fakeConsumer{}, decoder.NewOrderPaidDecoder(), producer, time.Millisecond)

		err := s.OrderPaidHandler(context.Background(), orderPaidMessage(t, "order-1"))

		require.NoError(t, err)
		require.Len(t, producer.events, 1)
		assert.Equal(t, "order-1", producer.events[0].OrderUUID)
		assert.Equal(t, "user-1", producer.events[0].UserUUID)
	})

	t.Run("нечитаемое сообщение не приводит к публикации", func(t *testing.T) {
		t.Parallel()

		producer := &fakeShipProducer{}
		s := NewService(&fakeConsumer{}, decoder.NewOrderPaidDecoder(), producer, time.Millisecond)

		msg := platformKafka.Message{Topic: "order.paid", Value: []byte("это не protobuf")}

		require.Error(t, s.OrderPaidHandler(context.Background(), msg))
		assert.Empty(t, producer.events)
	})

	t.Run("отмена контекста прерывает сборку", func(t *testing.T) {
		t.Parallel()

		producer := &fakeShipProducer{}
		s := NewService(&fakeConsumer{}, decoder.NewOrderPaidDecoder(), producer, time.Hour)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := s.OrderPaidHandler(ctx, orderPaidMessage(t, "order-2"))

		require.ErrorIs(t, err, context.Canceled)
		assert.Empty(t, producer.events)
	})

	t.Run("ошибка публикации возвращается наверх", func(t *testing.T) {
		t.Parallel()

		producer := &fakeShipProducer{err: errors.New("kafka down")}
		s := NewService(&fakeConsumer{}, decoder.NewOrderPaidDecoder(), producer, time.Millisecond)

		require.Error(t, s.OrderPaidHandler(context.Background(), orderPaidMessage(t, "order-3")))
	})
}
