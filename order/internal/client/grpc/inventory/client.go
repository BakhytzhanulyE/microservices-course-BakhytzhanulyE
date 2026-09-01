// Package inventory — gRPC-клиент каталога деталей.
package inventory

import (
	"context"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
	inventoryV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/inventory/v1"
)

type client struct {
	generatedClient inventoryV1.InventoryServiceClient
}

// NewClient оборачивает сгенерированный gRPC-клиент, чтобы наверх уезжала
// доменная модель, а не protobuf.
func NewClient(generatedClient inventoryV1.InventoryServiceClient) *client {
	return &client{generatedClient: generatedClient}
}

// ListParts возвращает детали по списку UUID.
func (c *client) ListParts(ctx context.Context, uuids []string) ([]model.Part, error) {
	resp, err := c.generatedClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: uuids,
		},
	})
	if err != nil {
		return nil, err
	}

	parts := make([]model.Part, 0, len(resp.GetParts()))
	for _, part := range resp.GetParts() {
		parts = append(parts, model.Part{
			UUID:  part.GetUuid(),
			Name:  part.GetName(),
			Price: part.GetPrice(),
		})
	}

	return parts, nil
}
