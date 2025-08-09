package inmemory_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
)

func (s *InMemoryRepositorySuite) TestListParts_Success() {
	// Вызываем метод без фильтра
	result, err := s.repository.ListParts(context.Background(), nil)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), *result, 2) // Должно быть 2 тестовых элемента

	// Проверяем, что все части присутствуют
	names := make([]string, 0, len(*result))
	for _, part := range *result {
		names = append(names, part.Name)
	}
	assert.Contains(s.T(), names, "Detail 1")
	assert.Contains(s.T(), names, "Detail 2")
}

func (s *InMemoryRepositorySuite) TestListParts_WithUUIDFilter() {
	// Тестируем фильтр по UUID
	testUUID := uuid.MustParse("d973e963-b7e6-4323-8f4e-4bfd5ab8e834")
	filter := &model.Filter{
		Uuids: []uuid.UUID{testUUID},
	}

	result, err := s.repository.ListParts(context.Background(), filter)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), *result, 1)
	assert.Equal(s.T(), "Detail 1", (*result)[0].Name)
	assert.Equal(s.T(), testUUID, (*result)[0].OrderUuid)
}

func (s *InMemoryRepositorySuite) TestListParts_WithNameFilter() {
	// Тестируем фильтр по имени
	filter := &model.Filter{
		Names: []string{"Detail 2"},
	}

	result, err := s.repository.ListParts(context.Background(), filter)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), *result, 1)
	assert.Equal(s.T(), "Detail 2", (*result)[0].Name)
	assert.Equal(s.T(), 200.0, (*result)[0].Price)
	assert.Equal(s.T(), int64(20), (*result)[0].StockQuantity)
}

func (s *InMemoryRepositorySuite) TestListParts_WithCategoryFilter() {
	// Тестируем фильтр по категории (оба тестовых элемента имеют категорию ENGINE)
	filter := &model.Filter{
		Categories: []model.Category{model.ENGINE},
	}

	result, err := s.repository.ListParts(context.Background(), filter)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), *result, 2) // Оба элемента должны пройти фильтр

	for _, part := range *result {
		assert.Equal(s.T(), model.ENGINE, part.Category)
	}
}

func (s *InMemoryRepositorySuite) TestListParts_WithManufacturerCountryFilter() {
	// Тестируем фильтр по стране производителя
	filter := &model.Filter{
		ManufacturerCountries: []string{"China"},
	}

	result, err := s.repository.ListParts(context.Background(), filter)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), *result, 1)
	assert.Equal(s.T(), "Detail 1", (*result)[0].Name)
	assert.Equal(s.T(), "China", (*result)[0].Manufacturer.Country)
}

func (s *InMemoryRepositorySuite) TestListParts_WithTagsFilter() {
	// Тестируем фильтр по тегам (оба элемента имеют теги tag1 и tag2)
	filter := &model.Filter{
		Tags: []string{"tag1"},
	}

	result, err := s.repository.ListParts(context.Background(), filter)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), *result, 2) // Оба элемента должны пройти фильтр

	for _, part := range *result {
		assert.Contains(s.T(), part.Tags, "tag1")
	}
}

func (s *InMemoryRepositorySuite) TestListParts_MultipleFilters() {
	// Тестируем комбинацию фильтров
	filter := &model.Filter{
		ManufacturerCountries: []string{"USA"},
		Tags:                  []string{"tag2"},
	}

	result, err := s.repository.ListParts(context.Background(), filter)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), *result, 1)
	assert.Equal(s.T(), "Detail 2", (*result)[0].Name)
	assert.Equal(s.T(), "USA", (*result)[0].Manufacturer.Country)
	assert.Contains(s.T(), (*result)[0].Tags, "tag2")
}

func (s *InMemoryRepositorySuite) TestListParts_NoMatches() {
	// Тестируем фильтр, который не должен ничего найти
	filter := &model.Filter{
		ManufacturerCountries: []string{"NonExistentCountry"},
	}

	result, err := s.repository.ListParts(context.Background(), filter)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), *result, 0)
}

func (s *InMemoryRepositorySuite) TestListParts_EmptyFilter() {
	// Тестируем пустой фильтр
	filter := &model.Filter{}

	result, err := s.repository.ListParts(context.Background(), filter)

	// Проверяем результат
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), *result, 2) // Пустой фильтр должен вернуть все элементы
}
