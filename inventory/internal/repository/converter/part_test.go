package converter

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/model"
)

func (s *ConverterSuite) TestToModelPart_Success() {
	// Подготовка
	partUUID := uuid.New()
	now := time.Now()

	repoPart := &repoModel.RepositoryPart{
		OrderUuid:     partUUID.String(),
		Name:          "Rocket Engine",
		Description:   "Powerful rocket engine",
		Price:         1500.75,
		StockQuantity: 5,
		Category:      1, // ENGINE
		Dimensions: repoModel.Dimensions{
			Length: 100.5,
			Width:  25.3,
			Height: 25.3,
			Weight: 150.8,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "SpaceX",
			Country: "USA",
			Website: "https://spacex.com",
		},
		Tags: []string{"engine", "rocket", "propulsion"},
		Metadata: map[string]repoModel.Value{
			"thrust": {
				Float64Value: 1000000.5,
			},
			"fuel_type": {
				StringValue: "RP-1",
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Выполнение
	result, err := ToModelPart(repoPart)

	// Проверка
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), partUUID, result.OrderUuid)
	assert.Equal(s.T(), "Rocket Engine", result.Name)
	assert.Equal(s.T(), "Powerful rocket engine", result.Description)
	assert.Equal(s.T(), 1500.75, result.Price)
	assert.Equal(s.T(), int64(5), result.StockQuantity)
	assert.Equal(s.T(), model.ENGINE, result.Category)

	// Проверка Dimensions
	assert.Equal(s.T(), 100.5, result.Dimensions.Length)
	assert.Equal(s.T(), 25.3, result.Dimensions.Width)
	assert.Equal(s.T(), 25.3, result.Dimensions.Height)
	assert.Equal(s.T(), 150.8, result.Dimensions.Weight)

	// Проверка Manufacturer
	assert.Equal(s.T(), "SpaceX", result.Manufacturer.Name)
	assert.Equal(s.T(), "USA", result.Manufacturer.Country)
	assert.Equal(s.T(), "https://spacex.com", result.Manufacturer.Website)

	// Проверка Tags
	assert.Equal(s.T(), []string{"engine", "rocket", "propulsion"}, result.Tags)

	// Проверка Metadata
	assert.Len(s.T(), result.Metadata, 2)
	assert.Equal(s.T(), float64(1000000.5), result.Metadata["thrust"].Float64Value)
	assert.Equal(s.T(), "RP-1", result.Metadata["fuel_type"].StringValue)

	// Проверка временных меток
	assert.Equal(s.T(), now, result.CreatedAt)
	assert.Equal(s.T(), now, result.UpdatedAt)
}

func (s *ConverterSuite) TestToModelPart_InvalidUUID() {
	// Подготовка
	repoPart := &repoModel.RepositoryPart{
		OrderUuid:     "invalid-uuid",
		Name:          "Rocket Engine",
		Description:   "Powerful rocket engine",
		Price:         1500.75,
		StockQuantity: 5,
		Category:      1,
		Dimensions:    repoModel.Dimensions{},
		Manufacturer:  repoModel.Manufacturer{},
		Tags:          []string{},
		Metadata:      map[string]repoModel.Value{},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Выполнение
	result, err := ToModelPart(repoPart)

	// Проверка
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), model.ErrConvertFromRepo, err)
}

func (s *ConverterSuite) TestToModelPart_EmptyUUID() {
	// Подготовка
	repoPart := &repoModel.RepositoryPart{
		OrderUuid:     "",
		Name:          "Rocket Engine",
		Description:   "Powerful rocket engine",
		Price:         1500.75,
		StockQuantity: 5,
		Category:      1,
		Dimensions:    repoModel.Dimensions{},
		Manufacturer:  repoModel.Manufacturer{},
		Tags:          []string{},
		Metadata:      map[string]repoModel.Value{},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Выполнение
	result, err := ToModelPart(repoPart)

	// Проверка
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), model.ErrConvertFromRepo, err)
}

func (s *ConverterSuite) TestToModelPart_AllCategories() {
	// Подготовка
	partUUID := uuid.New()
	now := time.Now()

	testCases := []struct {
		category   int
		expected   model.Category
		shouldPass bool
	}{
		{1, model.ENGINE, true},
		{2, model.FUEL, true},
		{3, model.PORTHOLE, true},
		{4, model.WING, true},
		{0, model.UNKNOWN, true},
		{999, model.UNKNOWN, true}, // неизвестная категория
	}

	for _, tc := range testCases {
		repoPart := &repoModel.RepositoryPart{
			OrderUuid:     partUUID.String(),
			Name:          "Test Part",
			Description:   "Test Description",
			Price:         100.0,
			StockQuantity: 1,
			Category:      tc.category,
			Dimensions:    repoModel.Dimensions{},
			Manufacturer:  repoModel.Manufacturer{},
			Tags:          []string{},
			Metadata:      map[string]repoModel.Value{},
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		// Выполнение
		result, err := ToModelPart(repoPart)

		// Проверка
		if tc.shouldPass {
			assert.NoError(s.T(), err)
			assert.NotNil(s.T(), result)
			assert.Equal(s.T(), tc.expected, result.Category)
		} else {
			assert.Error(s.T(), err)
			assert.Nil(s.T(), result)
		}
	}
}

func (s *ConverterSuite) TestToModelPart_AllMetadataTypes() {
	// Подготовка
	partUUID := uuid.New()
	now := time.Now()

	repoPart := &repoModel.RepositoryPart{
		OrderUuid:     partUUID.String(),
		Name:          "Test Part",
		Description:   "Test Description",
		Price:         100.0,
		StockQuantity: 1,
		Category:      1,
		Dimensions:    repoModel.Dimensions{},
		Manufacturer:  repoModel.Manufacturer{},
		Tags:          []string{},
		Metadata: map[string]repoModel.Value{
			"string_val": {
				StringValue: "test string",
			},
			"int_val": {
				Int64Value: 42,
			},
			"float_val": {
				Float64Value: 3.14,
			},
			"bool_val": {
				BoolValue: true,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Выполнение
	result, err := ToModelPart(repoPart)

	// Проверка
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)

	// Проверка всех типов метаданных
	assert.Equal(s.T(), "test string", result.Metadata["string_val"].StringValue)
	assert.Equal(s.T(), int64(42), result.Metadata["int_val"].Int64Value)
	assert.Equal(s.T(), 3.14, result.Metadata["float_val"].Float64Value)
	assert.Equal(s.T(), true, result.Metadata["bool_val"].BoolValue)
}

func (s *ConverterSuite) TestToModelPart_EmptyMetadata() {
	// Подготовка
	partUUID := uuid.New()
	now := time.Now()

	repoPart := &repoModel.RepositoryPart{
		OrderUuid:     partUUID.String(),
		Name:          "Test Part",
		Description:   "Test Description",
		Price:         100.0,
		StockQuantity: 1,
		Category:      1,
		Dimensions:    repoModel.Dimensions{},
		Manufacturer:  repoModel.Manufacturer{},
		Tags:          []string{},
		Metadata:      map[string]repoModel.Value{},
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Выполнение
	result, err := ToModelPart(repoPart)

	// Проверка
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.Metadata)
}

func (s *ConverterSuite) TestToModelPart_ZeroValues() {
	// Подготовка
	partUUID := uuid.New()
	now := time.Now()

	repoPart := &repoModel.RepositoryPart{
		OrderUuid:     partUUID.String(),
		Name:          "",
		Description:   "",
		Price:         0,
		StockQuantity: 0,
		Category:      0,
		Dimensions: repoModel.Dimensions{
			Length: 0,
			Width:  0,
			Height: 0,
			Weight: 0,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "",
			Country: "",
			Website: "",
		},
		Tags:      []string{},
		Metadata:  map[string]repoModel.Value{},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Выполнение
	result, err := ToModelPart(repoPart)

	// Проверка
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), partUUID, result.OrderUuid)
	assert.Equal(s.T(), "", result.Name)
	assert.Equal(s.T(), "", result.Description)
	assert.Equal(s.T(), float64(0), result.Price)
	assert.Equal(s.T(), int64(0), result.StockQuantity)
	assert.Equal(s.T(), model.UNKNOWN, result.Category)
	assert.Equal(s.T(), float64(0), result.Dimensions.Length)
	assert.Equal(s.T(), "", result.Manufacturer.Name)
	assert.Empty(s.T(), result.Tags)
	assert.Empty(s.T(), result.Metadata)
}
