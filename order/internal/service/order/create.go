package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
)

// Create проверяет детали в каталоге, считает сумму и создаёт заказ в статусе PENDING_PAYMENT.
func (s *service) Create(ctx context.Context, params model.CreateOrderParams) (model.Order, error) {
	if len(params.PartUUIDs) == 0 {
		return model.Order{}, model.ErrEmptyPartList
	}

	// Дубли в запросе схлопываем: цену за одну и ту же деталь считаем один раз,
	// иначе повтор UUID молча удвоил бы сумму заказа.
	uniqueUUIDs := uniqueStrings(params.PartUUIDs)

	parts, err := s.inventoryClient.ListParts(ctx, uniqueUUIDs)
	if err != nil {
		logger.Error(ctx, "Каталог деталей недоступен", zap.Error(err))
		return model.Order{}, model.ErrInventoryUnavailable
	}

	// Каталог возвращает только найденное, поэтому расхождение размеров означает,
	// что какой-то детали из заказа не существует.
	if len(parts) != len(uniqueUUIDs) {
		return model.Order{}, model.ErrPartsNotFound
	}

	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.Price
	}

	order := model.Order{
		UUID:       uuid.NewString(),
		UserUUID:   params.UserUUID,
		PartUUIDs:  uniqueUUIDs,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  time.Now().UTC(),
	}

	if err = s.orderRepository.Create(ctx, order); err != nil {
		return model.Order{}, err
	}

	logger.Info(ctx, "🧾 Заказ создан",
		zap.String("order_uuid", order.UUID),
		zap.String("user_uuid", order.UserUUID),
		zap.Float64("total_price", order.TotalPrice),
	)

	return order, nil
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))

	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}

		seen[value] = struct{}{}
		result = append(result, value)
	}

	return result
}
