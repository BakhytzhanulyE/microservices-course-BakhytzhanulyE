package payment

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/payment/internal/model"
)

func TestPayOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		method  model.PaymentMethod
		wantErr error
	}{
		{name: "карта", method: model.PaymentMethodCard},
		{name: "СБП", method: model.PaymentMethodSBP},
		{name: "кредитка", method: model.PaymentMethodCreditCard},
		{name: "деньги инвестора", method: model.PaymentMethodInvestorMoney},
		{name: "способ не указан", method: model.PaymentMethodUnspecified, wantErr: model.ErrUnsupportedPaymentMethod},
		{name: "неизвестный способ", method: model.PaymentMethod(42), wantErr: model.ErrUnsupportedPaymentMethod},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NewService()

			transactionUUID, err := s.PayOrder(context.Background(), model.Payment{
				OrderUUID:     "order-1",
				UserUUID:      "user-1",
				PaymentMethod: tt.method,
			})

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, transactionUUID)

				return
			}

			require.NoError(t, err)

			_, parseErr := uuid.Parse(transactionUUID)
			assert.NoError(t, parseErr, "сервис должен возвращать валидный UUID транзакции")
		})
	}
}

func TestPaymentMethodString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "CARD", model.PaymentMethodCard.String())
	assert.Equal(t, "INVESTOR_MONEY", model.PaymentMethodInvestorMoney.String())
	assert.Equal(t, "UNSPECIFIED", model.PaymentMethod(99).String())
}
