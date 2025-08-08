package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

// InsertTestSighting — вставляет тестовое наблюдение НЛО в коллекцию Mongo и возвращает его UUID
func (env *TestEnvironment) InsertTestSighting(ctx context.Context) (string, error) {
	sightingUUID := gofakeit.UUID()
	now := time.Now()

	sightingDoc := bson.M{
		"_id": sightingUUID,
		"info": bson.M{
			"observed_at":      primitive.NewDateTimeFromTime(now.Add(-time.Hour)),
			"location":         gofakeit.City() + ", " + gofakeit.Country(),
			"description":      gofakeit.Sentence(10),
			"color":            gofakeit.Color(),
			"sound":            gofakeit.Bool(),
			"duration_seconds": gofakeit.Number(10, 3600),
		},
		"created_at": primitive.NewDateTimeFromTime(now),
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "ufo-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(collectionName).InsertOne(ctx, sightingDoc)
	if err != nil {
		return "", err
	}

	return sightingUUID, nil
}

// InsertTestSightingWithData — вставляет тестовое наблюдение НЛО с заданными данными
func (env *TestEnvironment) InsertTestSightingWithData(ctx context.Context, info *inventoryV1.SightingInfo) (string, error) {
	sightingUUID := gofakeit.UUID()
	now := time.Now()

	observedAt := info.GetObservedAt().AsTime()

	sightingDoc := bson.M{
		"_id": sightingUUID,
		"info": bson.M{
			"observed_at":      primitive.NewDateTimeFromTime(observedAt),
			"location":         info.GetLocation(),
			"description":      info.GetDescription(),
			"color":            info.GetColor().GetValue(),
			"sound":            info.GetSound().GetValue(),
			"duration_seconds": info.GetDurationSeconds().GetValue(),
		},
		"created_at": primitive.NewDateTimeFromTime(now),
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "ufo-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(collectionName).InsertOne(ctx, sightingDoc)
	if err != nil {
		return "", err
	}

	return sightingUUID, nil
}

// GetTestSightingInfo — возвращает тестовую информацию о наблюдении НЛО
func (env *TestEnvironment) GetTestSightingInfo() *inventoryV1.SightingInfo {
	return &inventoryV1.SightingInfo{
		ObservedAt:      timestamppb.New(time.Now().Add(-2 * time.Hour)),
		Location:        "Москва, Красная площадь",
		Description:     "Яркий светящийся объект треугольной формы",
		Color:           wrapperspb.String("зеленый"),
		Sound:           wrapperspb.Bool(false),
		DurationSeconds: wrapperspb.Int32(120),
	}
}

// GetUpdatedSightingInfo — возвращает обновленную информацию о наблюдении НЛО
func (env *TestEnvironment) GetUpdatedSightingInfo() *inventoryV1.SightingUpdateInfo {
	return &inventoryV1.SightingUpdateInfo{
		Location:        wrapperspb.String("Санкт-Петербург, Дворцовая площадь"),
		Description:     wrapperspb.String("Обновленное описание: объект изменил форму"),
		Color:           wrapperspb.String("синий"),
		DurationSeconds: wrapperspb.Int32(180),
	}
}

// ClearSightingsCollection — удаляет все записи из коллекции sightings
func (env *TestEnvironment) ClearSightingsCollection(ctx context.Context) error {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "ufo-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(collectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
