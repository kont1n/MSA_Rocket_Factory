//go:build integration

package mongo_test

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/model"
)

func (s *MongoRepositorySuite) TestListParts_Success() {
	// Подготавливаем тестовые данные
	testParts := []*repoModel.RepositoryPart{
		{
			OrderUuid:     uuid.New().String(),
			Name:          "Rocket Engine",
			Description:   "Main rocket engine",
			Price:         20000.0,
			StockQuantity: 3,
			Category:      int(model.ENGINE),
			Dimensions: repoModel.Dimensions{
				Length: 150.0,
				Width:  60.0,
				Height: 80.0,
				Weight: 800.0,
			},
			Manufacturer: repoModel.Manufacturer{
				Name:    "SpaceX",
				Country: "USA",
			},
			Tags:      []string{"rocket", "engine", "main"},
			Metadata:  map[string]repoModel.Value{"thrust": {StringValue: "1000kN"}},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			OrderUuid:     uuid.New().String(),
			Name:          "Fuel Tank",
			Description:   "Large fuel storage tank",
			Price:         5000.0,
			StockQuantity: 10,
			Category:      int(model.FUEL),
			Dimensions: repoModel.Dimensions{
				Length: 300.0,
				Width:  100.0,
				Height: 100.0,
				Weight: 200.0,
			},
			Manufacturer: repoModel.Manufacturer{
				Name:    "Boeing",
				Country: "USA",
			},
			Tags:      []string{"fuel", "tank", "storage"},
			Metadata:  map[string]repoModel.Value{"capacity": {StringValue: "10000L"}},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Вставляем тестовые данные
	collection := s.db.Collection("parts")
	for _, part := range testParts {
		_, err := collection.InsertOne(context.Background(), part)
		s.Require().NoError(err)
	}

	// Вызываем метод без фильтра
	result, err := s.repository.ListParts(context.Background(), nil)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), *result, 2)

	// Проверяем, что все части присутствуют
	names := make([]string, 0, len(*result))
	for _, part := range *result {
		names = append(names, part.Name)
	}
	assert.Contains(s.T(), names, "Rocket Engine")
	assert.Contains(s.T(), names, "Fuel Tank")
}

func (s *MongoRepositorySuite) TestListParts_WithFilter() {
	// Подготавливаем тестовые данные
	engineUUID := uuid.New()
	testParts := []*repoModel.RepositoryPart{
		{
			OrderUuid:     engineUUID.String(),
			Name:          "Rocket Engine",
			Description:   "Main rocket engine",
			Price:         20000.0,
			StockQuantity: 3,
			Category:      int(model.ENGINE),
			Manufacturer: repoModel.Manufacturer{
				Name:    "SpaceX",
				Country: "USA",
			},
			Tags:      []string{"rocket", "engine"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			OrderUuid:     uuid.New().String(),
			Name:          "Fuel Tank",
			Description:   "Large fuel storage tank",
			Price:         5000.0,
			StockQuantity: 10,
			Category:      int(model.FUEL),
			Manufacturer: repoModel.Manufacturer{
				Name:    "Boeing",
				Country: "USA",
			},
			Tags:      []string{"fuel", "tank"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Вставляем тестовые данные
	collection := s.db.Collection("parts")
	for _, part := range testParts {
		_, err := collection.InsertOne(context.Background(), part)
		s.Require().NoError(err)
	}

	// Тестируем фильтр по UUID
	filter := &model.Filter{
		Uuids: []uuid.UUID{engineUUID},
	}

	result, err := s.repository.ListParts(context.Background(), filter)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), *result, 1)
	assert.Equal(s.T(), "Rocket Engine", (*result)[0].Name)
	assert.Equal(s.T(), engineUUID, (*result)[0].OrderUuid)
}

func (s *MongoRepositorySuite) TestListParts_WithCategoryFilter() {
	// Подготавливаем тестовые данные с разными категориями
	testParts := []*repoModel.RepositoryPart{
		{
			OrderUuid:     uuid.New().String(),
			Name:          "Rocket Engine",
			Category:      int(model.ENGINE),
			Price:         20000.0,
			StockQuantity: 3,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			OrderUuid:     uuid.New().String(),
			Name:          "Fuel Tank",
			Category:      int(model.FUEL),
			Price:         5000.0,
			StockQuantity: 10,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			OrderUuid:     uuid.New().String(),
			Name:          "Navigation System",
			Category:      int(model.WING),
			Price:         15000.0,
			StockQuantity: 2,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	// Вставляем тестовые данные
	collection := s.db.Collection("parts")
	for _, part := range testParts {
		_, err := collection.InsertOne(context.Background(), part)
		s.Require().NoError(err)
	}

	// Тестируем фильтр по категории
	filter := &model.Filter{
		Categories: []model.Category{model.ENGINE},
	}

	result, err := s.repository.ListParts(context.Background(), filter)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), *result, 1)
	assert.Equal(s.T(), "Rocket Engine", (*result)[0].Name)
	assert.Equal(s.T(), model.ENGINE, (*result)[0].Category)
}

func (s *MongoRepositorySuite) TestListParts_EmptyResult() {
	// Не вставляем никаких данных

	// Вызываем метод
	result, err := s.repository.ListParts(context.Background(), nil)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), *result, 0)
}

func (s *MongoRepositorySuite) TestListParts_WithTagsFilter() {
	// Подготавливаем тестовые данные
	testParts := []*repoModel.RepositoryPart{
		{
			OrderUuid:     uuid.New().String(),
			Name:          "Rocket Engine",
			Category:      int(model.ENGINE),
			Tags:          []string{"rocket", "engine", "propulsion"},
			Price:         20000.0,
			StockQuantity: 3,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			OrderUuid:     uuid.New().String(),
			Name:          "Fuel Tank",
			Category:      int(model.FUEL),
			Tags:          []string{"fuel", "tank", "storage"},
			Price:         5000.0,
			StockQuantity: 10,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	// Вставляем тестовые данные
	collection := s.db.Collection("parts")
	for _, part := range testParts {
		_, err := collection.InsertOne(context.Background(), part)
		s.Require().NoError(err)
	}

	// Тестируем фильтр по тегам
	filter := &model.Filter{
		Tags: []string{"engine"},
	}

	result, err := s.repository.ListParts(context.Background(), filter)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), *result, 1)
	assert.Equal(s.T(), "Rocket Engine", (*result)[0].Name)
	assert.Contains(s.T(), (*result)[0].Tags, "engine")
}
