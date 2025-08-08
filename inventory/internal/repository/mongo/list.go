package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
	repoConverter "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/converter"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/model"
)

func (r *repository) ListParts(ctx context.Context, filter *model.Filter) (*[]model.Part, error) {
	collection := r.db.Collection(partsCollection)

	// Создаем фильтр для MongoDB
	mongoFilter := buildMongoFilter(filter)

	// Выполняем запрос
	cursor, err := collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("failed to close cursor: %v", err)
		}
	}(cursor, ctx)

	// Читаем результаты
	var repoParts []repoModel.RepositoryPart
	if err = cursor.All(ctx, &repoParts); err != nil {
		return nil, err
	}

	// Преобразуем в основные модели
	parts := make([]model.Part, 0, len(repoParts))
	for _, repoPart := range repoParts {
		part, err := repoConverter.ToModelPart(&repoPart)
		if err != nil {
			return nil, err
		}
		parts = append(parts, *part)
	}

	return &parts, nil
}

// buildMongoFilter создание фильтра для MongoDB на основе модели фильтра
func buildMongoFilter(filter *model.Filter) bson.M {
	if filter == nil {
		return bson.M{}
	}

	mongoFilter := bson.M{}

	// Фильтр по UUID
	if len(filter.Uuids) > 0 {
		uuidStrings := make([]string, len(filter.Uuids))
		for i, uuid := range filter.Uuids {
			uuidStrings[i] = uuid.String()
		}
		mongoFilter["order_uuid"] = bson.M{"$in": uuidStrings}
	}

	// Фильтр по именам
	if len(filter.Names) > 0 {
		mongoFilter["name"] = bson.M{"$in": filter.Names}
	}

	// Фильтр по категориям
	if len(filter.Categories) > 0 {
		categoryInts := make([]int, len(filter.Categories))
		for i, category := range filter.Categories {
			categoryInts[i] = int(category)
		}
		mongoFilter["category"] = bson.M{"$in": categoryInts}
	}

	// Фильтр по странам производителей
	if len(filter.ManufacturerCountries) > 0 {
		mongoFilter["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}

	// Фильтр по тегам (должен содержать хотя бы один из указанных тегов)
	if len(filter.Tags) > 0 {
		mongoFilter["tags"] = bson.M{"$in": filter.Tags}
	}

	return mongoFilter
}
