// Package converter переводит детали между доменной моделью и protobuf.
package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/model"
	inventoryV1 "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/shared/pkg/proto/inventory/v1"
)

// PartToProto переводит доменную деталь в protobuf.
func PartToProto(part model.Part) *inventoryV1.Part {
	var updatedAt *timestamppb.Timestamp
	if part.UpdatedAt != nil {
		updatedAt = timestamppb.New(*part.UpdatedAt)
	}

	return &inventoryV1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      inventoryV1.Category(part.Category),
		Dimensions:    dimensionsToProto(part.Dimensions),
		Manufacturer:  manufacturerToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      metadataToProto(part.Metadata),
		CreatedAt:     timestamppb.New(part.CreatedAt),
		UpdatedAt:     updatedAt,
	}
}

// PartsToProto переводит список деталей в protobuf.
func PartsToProto(parts []model.Part) []*inventoryV1.Part {
	result := make([]*inventoryV1.Part, 0, len(parts))
	for _, part := range parts {
		result = append(result, PartToProto(part))
	}

	return result
}

// FilterToModel переводит protobuf-фильтр в доменный. nil означает «без фильтра».
func FilterToModel(filter *inventoryV1.PartsFilter) model.PartsFilter {
	if filter == nil {
		return model.PartsFilter{}
	}

	categories := make([]model.Category, 0, len(filter.GetCategories()))
	for _, category := range filter.GetCategories() {
		categories = append(categories, model.Category(category))
	}

	return model.PartsFilter{
		UUIDs:                 filter.GetUuids(),
		Names:                 filter.GetNames(),
		Categories:            categories,
		ManufacturerCountries: filter.GetManufacturerCountries(),
		Tags:                  filter.GetTags(),
	}
}

func dimensionsToProto(d *model.Dimensions) *inventoryV1.Dimensions {
	if d == nil {
		return nil
	}

	return &inventoryV1.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func manufacturerToProto(m *model.Manufacturer) *inventoryV1.Manufacturer {
	if m == nil {
		return nil
	}

	return &inventoryV1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func metadataToProto(meta map[string]model.Value) map[string]*inventoryV1.Value {
	if meta == nil {
		return nil
	}

	result := make(map[string]*inventoryV1.Value, len(meta))
	for key, value := range meta {
		result[key] = valueToProto(value)
	}

	return result
}

func valueToProto(value model.Value) *inventoryV1.Value {
	switch {
	case value.StringValue != nil:
		return &inventoryV1.Value{Kind: &inventoryV1.Value_StringValue{StringValue: *value.StringValue}}
	case value.Int64Value != nil:
		return &inventoryV1.Value{Kind: &inventoryV1.Value_Int64Value{Int64Value: *value.Int64Value}}
	case value.DoubleValue != nil:
		return &inventoryV1.Value{Kind: &inventoryV1.Value_DoubleValue{DoubleValue: *value.DoubleValue}}
	case value.BoolValue != nil:
		return &inventoryV1.Value{Kind: &inventoryV1.Value_BoolValue{BoolValue: *value.BoolValue}}
	case value.StringList != nil:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_StringList{
				StringList: &inventoryV1.StringList{Values: value.StringList},
			},
		}
	default:
		return &inventoryV1.Value{}
	}
}
