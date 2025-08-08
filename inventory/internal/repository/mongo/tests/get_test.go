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

func (s *MongoRepositorySuite) TestGetPart_Success() {
	// Подготавливаем тестовые данные
	partUUID := uuid.New()
	testPart := &repoModel.RepositoryPart{
		OrderUuid:     partUUID.String(),
		Name:          "Test Engine",
		Description:   "High performance rocket engine",
		Price:         15000.50,
		StockQuantity: 5,
		Category:      int(model.ENGINE),
		Dimensions: repoModel.Dimensions{
			Length: 100.0,
			Width:  50.0,
			Height: 75.0,
			Weight: 500.0,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "SpaceX",
			Country: "USA",
		},
		Tags: []string{"rocket", "engine", "high-performance"},
		Metadata: map[string]repoModel.Value{
			"thrust": {StringValue: "500kN"},
			"fuel":   {StringValue: "RP-1"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Вставляем тестовые данные в MongoDB
	collection := s.db.Collection("parts")
	_, err := collection.InsertOne(context.Background(), testPart)
	s.Require().NoError(err)

	// Вызываем метод репозитория
	result, err := s.repository.GetPart(context.Background(), partUUID)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), partUUID, result.OrderUuid)
	assert.Equal(s.T(), "Test Engine", result.Name)
	assert.Equal(s.T(), "High performance rocket engine", result.Description)
	assert.Equal(s.T(), 15000.50, result.Price)
	assert.Equal(s.T(), int64(5), result.StockQuantity)
	assert.Equal(s.T(), model.ENGINE, result.Category)
	assert.Equal(s.T(), "SpaceX", result.Manufacturer.Name)
	assert.Equal(s.T(), "USA", result.Manufacturer.Country)
	assert.Contains(s.T(), result.Tags, "rocket")
	assert.Contains(s.T(), result.Tags, "engine")
}

func (s *MongoRepositorySuite) TestGetPart_NotFound() {
	// Используем несуществующий UUID
	nonExistentUUID := uuid.New()

	// Вызываем метод репозитория
	result, err := s.repository.GetPart(context.Background(), nonExistentUUID)

	// Проверяем результат
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.ErrorIs(s.T(), err, model.ErrPartNotFound)
}

func (s *MongoRepositorySuite) TestGetPart_DatabaseError() {
	// Закрываем соединение с базой данных для симуляции ошибки
	err := s.client.Disconnect(context.Background())
	s.Require().NoError(err)

	partUUID := uuid.New()

	// Вызываем метод репозитория
	result, err := s.repository.GetPart(context.Background(), partUUID)

	// Проверяем результат
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.NotErrorIs(s.T(), err, model.ErrPartNotFound)

	// Восстанавливаем соединение для последующих тестов
	s.SetupSuite()
}
