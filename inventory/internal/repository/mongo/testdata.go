package mongo

import (
	"context"
	"log"
	"time"

	repoModel "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/model"
)

// AddTestData Добавление тестовых данных в MongoDB коллекцию
func (r *repository) AddTestData(ctx context.Context) error {
	log.Printf("Добавление тестовых данных для инвентаря в MongoDB")

	collection := r.db.Collection(partsCollection)

	// Проверяем, есть ли уже данные в коллекции
	count, err := collection.CountDocuments(ctx, map[string]interface{}{})
	if err != nil {
		return err
	}

	if count > 0 {
		log.Printf("Тестовые данные уже существуют в коллекции, пропускаем добавление")
		return nil
	}

	// Создаем тестовые данные
	testParts := []repoModel.RepositoryPart{
		{
			OrderUuid:     "d973e963-b7e6-4323-8f4e-4bfd5ab8e834",
			Name:          "Ракетный двигатель RD-180",
			Description:   "Мощный ракетный двигатель для тяжелых носителей",
			Price:         15000000.50,
			StockQuantity: 5,
			Category:      1, // ENGINE
			Dimensions: repoModel.Dimensions{
				Length: 355.6,
				Width:  297.2,
				Height: 297.2,
				Weight: 5480.0,
			},
			Manufacturer: repoModel.Manufacturer{
				Name:    "Energomash",
				Country: "Russia",
				Website: "https://www.energomash.ru",
			},
			Tags: []string{"rocket", "engine", "rd-180", "heavy-lift"},
			Metadata: map[string]repoModel.Value{
				"thrust": {
					Float64Value: 3830000.0,
				},
				"fuel_type": {
					StringValue: "RP-1/LOX",
				},
				"reusable": {
					BoolValue: false,
				},
				"test_fires": {
					Int64Value: 150,
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			OrderUuid:     "d973e963-b7e6-4323-8f4e-4bfd5ab8e835",
			Name:          "Система управления Navigation-1",
			Description:   "Навигационная система для космических аппаратов",
			Price:         750000.00,
			StockQuantity: 12,
			Category:      2, // ELECTRONICS
			Dimensions: repoModel.Dimensions{
				Length: 50.0,
				Width:  40.0,
				Height: 15.0,
				Weight: 8.5,
			},
			Manufacturer: repoModel.Manufacturer{
				Name:    "SpaceX",
				Country: "USA",
				Website: "https://www.spacex.com",
			},
			Tags: []string{"navigation", "electronics", "gps", "space"},
			Metadata: map[string]repoModel.Value{
				"accuracy": {
					StringValue: "sub-meter",
				},
				"frequency": {
					Float64Value: 1575.42,
				},
				"channels": {
					Int64Value: 32,
				},
				"space_qualified": {
					BoolValue: true,
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			OrderUuid:     "d973e963-b7e6-4323-8f4e-4bfd5ab8e836",
			Name:          "Топливный бак Falcon-Tank-9",
			Description:   "Алюминиевый топливный бак для среднего класса ракет",
			Price:         2500000.75,
			StockQuantity: 8,
			Category:      3, // STRUCTURE
			Dimensions: repoModel.Dimensions{
				Length: 1200.0,
				Width:  366.0,
				Height: 366.0,
				Weight: 25000.0,
			},
			Manufacturer: repoModel.Manufacturer{
				Name:    "Boeing",
				Country: "USA",
				Website: "https://www.boeing.com/space",
			},
			Tags: []string{"fuel", "tank", "structure", "aluminum"},
			Metadata: map[string]repoModel.Value{
				"capacity": {
					Float64Value: 400000.0,
				},
				"material": {
					StringValue: "Al-Li 2195",
				},
				"pressure_rating": {
					Float64Value: 3.5,
				},
				"insulated": {
					BoolValue: true,
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Конвертируем в interface{} для вставки
	docs := make([]interface{}, len(testParts))
	for i, part := range testParts {
		docs[i] = part
	}

	// Вставляем тестовые данные
	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		log.Printf("Ошибка при добавлении тестовых данных: %v", err)
		return err
	}

	log.Printf("Успешно добавлено %d тестовых записей в коллекцию %s", len(result.InsertedIDs), partsCollection)
	return nil
}
