// Package payment — gRPC-клиент платёжного сервиса.
package payment

import (
	"context"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
	paymentV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/payment/v1"
)

// methodByName переводит способ оплаты из HTTP-запроса в enum protobuf.
var methodByName = map[string]paymentV1.PaymentMethod{
	"CARD":           paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	"SBP":            paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
	"CREDIT_CARD":    paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
	"INVESTOR_MONEY": paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
}

type client struct {
	generatedClient paymentV1.PaymentServiceClient
}

// NewClient оборачивает сгенерированный gRPC-клиент платежей.
func NewClient(generatedClient paymentV1.PaymentServiceClient) *client {
	return &client{generatedClient: generatedClient}
}

// PayOrder проводит оплату и возвращает UUID транзакции.
func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error) {
	method, ok := methodByName[paymentMethod]
	if !ok {
		return "", model.ErrUnsupportedPaymentMethod
	}

	resp, err := c.generatedClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: method,
	})
	if err != nil {
		return "", err
	}

	return resp.GetTransactionUuid(), nil
}

// SupportedMethods возвращает список поддерживаемых способов оплаты —
// нужен HTTP-слою для валидации до похода в gRPC.
func SupportedMethods() []string {
	methods := make([]string, 0, len(methodByName))
	for name := range methodByName {
		methods = append(methods, name)
	}

	return methods
}

// IsSupportedMethod проверяет способ оплаты.
func IsSupportedMethod(name string) bool {
	_, ok := methodByName[name]
	return ok
}
