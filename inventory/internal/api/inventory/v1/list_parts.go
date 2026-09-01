package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/converter"
	inventoryV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/inventory/v1"
)

// ListParts возвращает детали, подходящие под фильтр.
func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := a.partService.List(ctx, converter.FilterToModel(req.GetFilter()))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list parts")
	}

	return &inventoryV1.ListPartsResponse{
		Parts: converter.PartsToProto(parts),
	}, nil
}
