// Package converter переводит детали между доменной моделью и моделью базы.
package converter

import (
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/model"
	repoModel "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/repository/model"
)

// PartToModel переводит деталь из модели базы в доменную.
func PartToModel(part repoModel.Part) model.Part {
	return model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.Category),
		Dimensions:    dimensionsToModel(part.Dimensions),
		Manufacturer:  manufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      metadataToModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

// PartToRepoModel переводит доменную деталь в модель базы.
func PartToRepoModel(part model.Part) repoModel.Part {
	return repoModel.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      int32(part.Category),
		Dimensions:    dimensionsToRepoModel(part.Dimensions),
		Manufacturer:  manufacturerToRepoModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      metadataToRepoModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func dimensionsToModel(d *repoModel.Dimensions) *model.Dimensions {
	if d == nil {
		return nil
	}

	return &model.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func dimensionsToRepoModel(d *model.Dimensions) *repoModel.Dimensions {
	if d == nil {
		return nil
	}

	return &repoModel.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func manufacturerToModel(m *repoModel.Manufacturer) *model.Manufacturer {
	if m == nil {
		return nil
	}

	return &model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func manufacturerToRepoModel(m *model.Manufacturer) *repoModel.Manufacturer {
	if m == nil {
		return nil
	}

	return &repoModel.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func metadataToModel(meta map[string]repoModel.Value) map[string]model.Value {
	if meta == nil {
		return nil
	}

	result := make(map[string]model.Value, len(meta))
	for k, v := range meta {
		result[k] = model.Value{
			StringValue: v.StringValue,
			Int64Value:  v.Int64Value,
			DoubleValue: v.DoubleValue,
			BoolValue:   v.BoolValue,
			StringList:  v.StringList,
		}
	}

	return result
}

func metadataToRepoModel(meta map[string]model.Value) map[string]repoModel.Value {
	if meta == nil {
		return nil
	}

	result := make(map[string]repoModel.Value, len(meta))
	for k, v := range meta {
		result[k] = repoModel.Value{
			StringValue: v.StringValue,
			Int64Value:  v.Int64Value,
			DoubleValue: v.DoubleValue,
			BoolValue:   v.BoolValue,
			StringList:  v.StringList,
		}
	}

	return result
}
