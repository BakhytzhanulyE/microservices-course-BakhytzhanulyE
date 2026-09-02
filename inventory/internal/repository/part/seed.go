package part

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/model"
	repoConverter "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/repository/converter"
)

// Seed наполняет каталог демо-данными, если он пуст.
// Нужен, чтобы сразу после docker compose up было что заказывать.
func (r *repository) Seed(ctx context.Context) error {
	count, err := r.collection.CountDocuments(ctx, map[string]any{})
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	now := time.Now().UTC()

	parts := []model.Part{
		{
			UUID:          uuid.NewString(),
			Name:          "Ионный двигатель IX-9",
			Description:   "Маршевый ионный двигатель для дальних перелётов",
			Price:         125000,
			StockQuantity: 12,
			Category:      model.CategoryEngine,
			Dimensions:    &model.Dimensions{Length: 320, Width: 180, Height: 180, Weight: 1450},
			Manufacturer:  &model.Manufacturer{Name: "Astra Dynamics", Country: "KZ", Website: "https://astra.example"},
			Tags:          []string{"engine", "ion", "long-range"},
			Metadata: map[string]model.Value{
				"thrust_kn":  {DoubleValue: lo.ToPtr(240.5)},
				"certified":  {BoolValue: lo.ToPtr(true)},
				"compatible": {StringList: []string{"X-1", "X-2"}},
			},
			CreatedAt: now,
		},
		{
			UUID:          uuid.NewString(),
			Name:          "Топливный блок H-200",
			Description:   "Криогенный водородный бак на 200 тонн",
			Price:         48000,
			StockQuantity: 40,
			Category:      model.CategoryFuel,
			Dimensions:    &model.Dimensions{Length: 600, Width: 240, Height: 240, Weight: 3200},
			Manufacturer:  &model.Manufacturer{Name: "Baikonur Fuels", Country: "KZ", Website: "https://fuels.example"},
			Tags:          []string{"fuel", "hydrogen"},
			Metadata: map[string]model.Value{
				"capacity_tons": {Int64Value: lo.ToPtr(int64(200))},
				"grade":         {StringValue: lo.ToPtr("A+")},
			},
			CreatedAt: now,
		},
		{
			UUID:          uuid.NewString(),
			Name:          "Иллюминатор Panorama-4",
			Description:   "Панорамный иллюминатор с тройным остеклением",
			Price:         15500,
			StockQuantity: 75,
			Category:      model.CategoryPorthole,
			Dimensions:    &model.Dimensions{Length: 90, Width: 90, Height: 12, Weight: 48},
			Manufacturer:  &model.Manufacturer{Name: "Orbital Glass", Country: "DE", Website: "https://glass.example"},
			Tags:          []string{"porthole", "panoramic"},
			CreatedAt:     now,
		},
		{
			UUID:          uuid.NewString(),
			Name:          "Крыло Delta-7",
			Description:   "Складное крыло для атмосферного манёвра",
			Price:         86000,
			StockQuantity: 20,
			Category:      model.CategoryWing,
			Dimensions:    &model.Dimensions{Length: 850, Width: 310, Height: 60, Weight: 970},
			Manufacturer:  &model.Manufacturer{Name: "Delta Aero", Country: "US", Website: "https://delta.example"},
			Tags:          []string{"wing", "foldable"},
			CreatedAt:     now,
		},
	}

	docs := make([]any, 0, len(parts))
	for _, part := range parts {
		docs = append(docs, repoConverter.PartToRepoModel(part))
	}

	_, err = r.collection.InsertMany(ctx, docs)

	return err
}
