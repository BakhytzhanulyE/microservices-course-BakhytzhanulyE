package part

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/model"
	repoConverter "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/repository/converter"
	repoModel "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/inventory/internal/repository/model"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
)

// List возвращает детали, подходящие под фильтр. Пустой фильтр отдаёт весь каталог.
func (r *repository) List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	cursor, err := r.collection.Find(ctx, buildFilter(filter))
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := cursor.Close(ctx); closeErr != nil {
			logger.Error(ctx, "Не удалось закрыть курсор MongoDB", zap.Error(closeErr))
		}
	}()

	var repoParts []repoModel.Part
	if err = cursor.All(ctx, &repoParts); err != nil {
		return nil, err
	}

	parts := make([]model.Part, 0, len(repoParts))
	for _, repoPart := range repoParts {
		parts = append(parts, repoConverter.PartToModel(repoPart))
	}

	return parts, nil
}

// buildFilter собирает bson-запрос: заполненные поля соединяются по «И»,
// значения внутри поля — по «ИЛИ» через $in.
func buildFilter(filter model.PartsFilter) bson.M {
	query := bson.M{}

	if len(filter.UUIDs) > 0 {
		query["_id"] = bson.M{"$in": filter.UUIDs}
	}

	if len(filter.Names) > 0 {
		query["name"] = bson.M{"$in": filter.Names}
	}

	if len(filter.Categories) > 0 {
		categories := make([]int32, 0, len(filter.Categories))
		for _, category := range filter.Categories {
			categories = append(categories, int32(category))
		}

		query["category"] = bson.M{"$in": categories}
	}

	if len(filter.ManufacturerCountries) > 0 {
		query["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}

	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}

	return query
}
