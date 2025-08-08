package inmemory_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
)

func (s *InMemoryRepositorySuite) TestGetPart_Success() {
	// Используем один из предустановленных UUID из тестовых данных
	testUUID := uuid.MustParse("d973e963-b7e6-4323-8f4e-4bfd5ab8e834")

	// Вызываем метод репозитория
	result, err := s.repository.GetPart(context.Background(), testUUID)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), testUUID, result.OrderUuid)
	assert.Equal(s.T(), "Detail 1", result.Name)
	assert.Equal(s.T(), "Detail 1 description", result.Description)
	assert.Equal(s.T(), 100.0, result.Price)
	assert.Equal(s.T(), int64(10), result.StockQuantity)
	assert.Equal(s.T(), model.ENGINE, result.Category)
	assert.Equal(s.T(), "China", result.Manufacturer.Country)
	assert.Equal(s.T(), "Details Fabric", result.Manufacturer.Name)
	assert.Contains(s.T(), result.Tags, "tag1")
	assert.Contains(s.T(), result.Tags, "tag2")
	assert.Equal(s.T(), 100.0, result.Dimensions.Length)
	assert.Equal(s.T(), 100.0, result.Dimensions.Width)
	assert.Equal(s.T(), 100.0, result.Dimensions.Height)
	assert.Equal(s.T(), 100.0, result.Dimensions.Weight)
}

func (s *InMemoryRepositorySuite) TestGetPart_NotFound() {
	// Используем несуществующий UUID
	nonExistentUUID := uuid.New()

	// Вызываем метод репозитория
	result, err := s.repository.GetPart(context.Background(), nonExistentUUID)

	// Проверяем результат
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.ErrorIs(s.T(), err, model.ErrPartNotFound)
}

func (s *InMemoryRepositorySuite) TestGetPart_ValidateTestData() {
	// Проверяем, что оба тестовых элемента доступны
	testUUIDs := []string{
		"d973e963-b7e6-4323-8f4e-4bfd5ab8e834",
		"d973e963-b7e6-4323-8f4e-4bfd5ab8e835",
	}

	expectedNames := []string{"Detail 1", "Detail 2"}
	expectedPrices := []float64{100.0, 200.0}
	expectedQuantities := []int64{10, 20}
	expectedCountries := []string{"China", "USA"}

	for i, uuidStr := range testUUIDs {
		testUUID := uuid.MustParse(uuidStr)
		result, err := s.repository.GetPart(context.Background(), testUUID)

		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), testUUID, result.OrderUuid)
		assert.Equal(s.T(), expectedNames[i], result.Name)
		assert.Equal(s.T(), expectedPrices[i], result.Price)
		assert.Equal(s.T(), expectedQuantities[i], result.StockQuantity)
		assert.Equal(s.T(), expectedCountries[i], result.Manufacturer.Country)
		assert.Equal(s.T(), model.ENGINE, result.Category)
	}
}
