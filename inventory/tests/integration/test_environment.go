//go:build integration

package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

// InsertTestPart — вставляет тестовую деталь ракеты в коллекцию Mongo и возвращает её UUID
func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUUID := gofakeit.UUID()
	now := time.Now()

	partDoc := bson.M{
		"_id":            primitive.NewObjectID(),
		"order_uuid":     partUUID,
		"name":           "Ракетный двигатель RD-180",
		"description":    "Мощный ракетный двигатель для тяжелых носителей",
		"price":          15000000.50,
		"stock_quantity": 5,
		"category":       1, // ENGINE
		"dimensions": bson.M{
			"length": 355.6,
			"width":  297.2,
			"height": 297.2,
			"weight": 5480.0,
		},
		"manufacturer": bson.M{
			"name":    "Energomash",
			"country": "Russia",
			"url":     "https://www.energomash.ru",
		},
		"tags": []string{"rocket", "engine", "rd-180", "heavy-lift"},
		"metadata": bson.M{
			"thrust": bson.M{
				"double_value": 3830000.0,
			},
			"fuel_type": bson.M{
				"string_value": "RP-1/LOX",
			},
		},
		"created_at": primitive.NewDateTimeFromTime(now),
		"updated_at": primitive.NewDateTimeFromTime(now),
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(collectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

// InsertTestPartWithData — вставляет тестовую деталь ракеты с заданными данными
func (env *TestEnvironment) InsertTestPartWithData(ctx context.Context, part *inventoryV1.Part) (string, error) {
	partUUID := gofakeit.UUID()
	now := time.Now()

	metadata := make(bson.M)
	for key, value := range part.GetMetadata() {
		switch v := value.GetKind().(type) {
		case *inventoryV1.Value_StringValue:
			metadata[key] = bson.M{"string_value": v.StringValue}
		case *inventoryV1.Value_Int64Value:
			metadata[key] = bson.M{"int64_value": v.Int64Value}
		case *inventoryV1.Value_DoubleValue:
			metadata[key] = bson.M{"double_value": v.DoubleValue}
		case *inventoryV1.Value_BoolValue:
			metadata[key] = bson.M{"bool_value": v.BoolValue}
		}
	}

	partDoc := bson.M{
		"_id":            primitive.NewObjectID(),
		"order_uuid":     partUUID,
		"name":           part.GetName(),
		"description":    part.GetDescription(),
		"price":          part.GetPrice(),
		"stock_quantity": part.GetStockQuantity(),
		"category":       int(part.GetCategory()),
		"dimensions": bson.M{
			"length": part.GetDimensions().GetLength(),
			"width":  part.GetDimensions().GetWidth(),
			"height": part.GetDimensions().GetHeight(),
			"weight": part.GetDimensions().GetWeight(),
		},
		"manufacturer": bson.M{
			"name":    part.GetManufacturer().GetName(),
			"country": part.GetManufacturer().GetCountry(),
			"url":     part.GetManufacturer().GetUrl(),
		},
		"tags":       part.GetTags(),
		"metadata":   metadata,
		"created_at": primitive.NewDateTimeFromTime(now),
		"updated_at": primitive.NewDateTimeFromTime(now),
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(collectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

// GetTestPart — возвращает тестовую информацию о детали ракеты
func (env *TestEnvironment) GetTestPart() *inventoryV1.Part {
	return &inventoryV1.Part{
		PartUuid:      gofakeit.UUID(),
		Name:          "Ракетный двигатель Merlin 1D",
		Description:   "Высокопроизводительный ракетный двигатель",
		Price:         12500000.00,
		StockQuantity: 8,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 300.0,
			Width:  150.0,
			Height: 150.0,
			Weight: 470.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "SpaceX",
			Country: "USA",
			Url:     "https://www.spacex.com",
		},
		Tags: []string{"rocket", "engine", "merlin", "reusable"},
		Metadata: map[string]*inventoryV1.Value{
			"thrust": {
				Kind: &inventoryV1.Value_DoubleValue{
					DoubleValue: 845000.0,
				},
			},
			"fuel_type": {
				Kind: &inventoryV1.Value_StringValue{
					StringValue: "RP-1/LOX",
				},
			},
			"reusable": {
				Kind: &inventoryV1.Value_BoolValue{
					BoolValue: true,
				},
			},
		},
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
	}
}

// GetTestFuelTankPart — возвращает тестовую информацию о топливном баке
func (env *TestEnvironment) GetTestFuelTankPart() *inventoryV1.Part {
	return &inventoryV1.Part{
		PartUuid:      gofakeit.UUID(),
		Name:          "Топливный бак Falcon-Tank-9",
		Description:   "Алюминиевый топливный бак для среднего класса ракет",
		Price:         2500000.75,
		StockQuantity: 3,
		Category:      inventoryV1.Category_CATEGORY_FUEL,
		Dimensions: &inventoryV1.Dimensions{
			Length: 1200.0,
			Width:  366.0,
			Height: 366.0,
			Weight: 25000.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "Boeing",
			Country: "USA",
			Url:     "https://www.boeing.com/space",
		},
		Tags: []string{"fuel", "tank", "structure", "aluminum"},
		Metadata: map[string]*inventoryV1.Value{
			"capacity": {
				Kind: &inventoryV1.Value_DoubleValue{
					DoubleValue: 400000.0,
				},
			},
			"material": {
				Kind: &inventoryV1.Value_StringValue{
					StringValue: "Al-Li 2195",
				},
			},
		},
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
	}
}

// InsertMultipleTestParts — вставляет несколько тестовых деталей для тестирования фильтрации
func (env *TestEnvironment) InsertMultipleTestParts(ctx context.Context) ([]string, error) {
	var partUUIDs []string
	now := time.Now()

	// Создаем UUID для каждой детали
	uuid1 := gofakeit.UUID()
	uuid2 := gofakeit.UUID()
	uuid3 := gofakeit.UUID()

	partUUIDs = append(partUUIDs, uuid1, uuid2, uuid3)

	testParts := []bson.M{
		{
			"_id":            primitive.NewObjectID(),
			"order_uuid":     uuid1,
			"name":           "Ракетный двигатель RD-180",
			"description":    "Мощный ракетный двигатель для тяжелых носителей",
			"price":          15000000.50,
			"stock_quantity": 5,
			"category":       1, // ENGINE
			"manufacturer": bson.M{
				"name":    "Energomash",
				"country": "Russia",
			},
			"tags":       []string{"rocket", "engine", "rd-180"},
			"created_at": primitive.NewDateTimeFromTime(now),
			"updated_at": primitive.NewDateTimeFromTime(now),
		},
		{
			"_id":            primitive.NewObjectID(),
			"order_uuid":     uuid2,
			"name":           "Топливный бак Falcon-Tank-9",
			"description":    "Алюминиевый топливный бак",
			"price":          2500000.75,
			"stock_quantity": 3,
			"category":       2, // FUEL
			"manufacturer": bson.M{
				"name":    "Boeing",
				"country": "USA",
			},
			"tags":       []string{"fuel", "tank", "aluminum"},
			"created_at": primitive.NewDateTimeFromTime(now),
			"updated_at": primitive.NewDateTimeFromTime(now),
		},
		{
			"_id":            primitive.NewObjectID(),
			"order_uuid":     uuid3,
			"name":           "Иллюминатор космический",
			"description":    "Прочный иллюминатор для космических кораблей",
			"price":          750000.00,
			"stock_quantity": 10,
			"category":       3, // PORTHOLE
			"manufacturer": bson.M{
				"name":    "Roscosmos",
				"country": "Russia",
			},
			"tags":       []string{"porthole", "window", "space"},
			"created_at": primitive.NewDateTimeFromTime(now),
			"updated_at": primitive.NewDateTimeFromTime(now),
		},
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	collection := env.Mongo.Client().Database(databaseName).Collection(collectionName)

	for _, part := range testParts {
		_, err := collection.InsertOne(ctx, part)
		if err != nil {
			return nil, err
		}
	}

	return partUUIDs, nil
}

// ClearPartsCollection — удаляет все записи из коллекции parts
func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(collectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
