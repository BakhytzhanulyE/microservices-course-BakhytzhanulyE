package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/payment/internal/model"
	paymentV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/payment/v1"
)

// PayOrder проводит оплату заказа.
func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	if req.GetOrderUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "order_uuid is required")
	}

	if req.GetUserUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_uuid is required")
	}

	transactionUUID, err := a.paymentService.PayOrder(ctx, model.Payment{
		OrderUUID:     req.GetOrderUuid(),
		UserUUID:      req.GetUserUuid(),
		PaymentMethod: model.PaymentMethod(req.GetPaymentMethod()),
	})
	if err != nil {
		if errors.Is(err, model.ErrUnsupportedPaymentMethod) {
			return nil, status.Error(codes.InvalidArgument, "unsupported payment method")
		}

		return nil, status.Error(codes.Internal, "failed to pay order")
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
