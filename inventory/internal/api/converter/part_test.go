package converter

import (
	"time"

	"github.com/google/uuid"
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
	"github.com/stretchr/testify/assert"
)

func (s *ConverterSuite) TestToProtoPart_EngineCategory() {
	// Подготовка
	partUUID := uuid.New()
	now := time.Now()
	part := &model.Part{
		OrderUuid:     partUUID,
		Name:          "Rocket Engine",
		Description:   "Powerful rocket engine",
		Price:         1500.75,
		StockQuantity: 5,
		Category:      model.ENGINE,
		Dimensions: model.Dimensions{
			Length: 100.5,
			Width:  25.3,
			Height: 25.3,
			Weight: 150.8,
		},
		Manufacturer: model.Manufacturer{
			Name:    "SpaceX",
			Country: "USA",
			Website: "https://spacex.com",
		},
		Tags: []string{"engine", "rocket", "propulsion"},
		Metadata: map[string]model.Value{
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
	result := ToProtoPart(part)

	// Проверка
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), partUUID.String(), result.PartUuid)
	assert.Equal(s.T(), "Rocket Engine", result.Name)
	assert.Equal(s.T(), "Powerful rocket engine", result.Description)
	assert.Equal(s.T(), 1500.75, result.Price)
	assert.Equal(s.T(), int64(5), result.StockQuantity)
	assert.Equal(s.T(), inventoryV1.Category_CATEGORY_ENGINE, result.Category)

	// Проверка Dimensions
	assert.Equal(s.T(), 100.5, result.Dimensions.Length)
	assert.Equal(s.T(), 25.3, result.Dimensions.Width)
	assert.Equal(s.T(), 25.3, result.Dimensions.Height)
	assert.Equal(s.T(), 150.8, result.Dimensions.Weight)

	// Проверка Manufacturer
	assert.Equal(s.T(), "SpaceX", result.Manufacturer.Name)
	assert.Equal(s.T(), "USA", result.Manufacturer.Country)
	assert.Equal(s.T(), "https://spacex.com", result.Manufacturer.Url)

	// Проверка Tags
	assert.Equal(s.T(), []string{"engine", "rocket", "propulsion"}, result.Tags)

	// Проверка Metadata
	assert.Len(s.T(), result.Metadata, 2)
	assert.Equal(s.T(), float64(1000000.5), result.Metadata["thrust"].GetDoubleValue())
	assert.Equal(s.T(), "RP-1", result.Metadata["fuel_type"].GetStringValue())

	// Проверка временных меток
	assert.NotNil(s.T(), result.CreatedAt)
	assert.NotNil(s.T(), result.UpdatedAt)
}

func (s *ConverterSuite) TestToProtoPart_FuelCategory() {
	part := &model.Part{
		OrderUuid:    uuid.New(),
		Name:         "Rocket Fuel",
		Category:     model.FUEL,
		Dimensions:   model.Dimensions{},
		Manufacturer: model.Manufacturer{},
		Tags:         []string{},
		Metadata:     map[string]model.Value{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	result := ToProtoPart(part)

	assert.Equal(s.T(), inventoryV1.Category_CATEGORY_FUEL, result.Category)
}

func (s *ConverterSuite) TestToProtoPart_PortholeCategory() {
	part := &model.Part{
		OrderUuid:    uuid.New(),
		Name:         "Space Porthole",
		Category:     model.PORTHOLE,
		Dimensions:   model.Dimensions{},
		Manufacturer: model.Manufacturer{},
		Tags:         []string{},
		Metadata:     map[string]model.Value{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	result := ToProtoPart(part)

	assert.Equal(s.T(), inventoryV1.Category_CATEGORY_PORTHOLE, result.Category)
}

func (s *ConverterSuite) TestToProtoPart_WingCategory() {
	part := &model.Part{
		OrderUuid:    uuid.New(),
		Name:         "Rocket Wing",
		Category:     model.WING,
		Dimensions:   model.Dimensions{},
		Manufacturer: model.Manufacturer{},
		Tags:         []string{},
		Metadata:     map[string]model.Value{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	result := ToProtoPart(part)

	assert.Equal(s.T(), inventoryV1.Category_CATEGORY_WING, result.Category)
}

func (s *ConverterSuite) TestToProtoPart_UnknownCategory() {
	part := &model.Part{
		OrderUuid:    uuid.New(),
		Name:         "Unknown Part",
		Category:     model.UNKNOWN,
		Dimensions:   model.Dimensions{},
		Manufacturer: model.Manufacturer{},
		Tags:         []string{},
		Metadata:     map[string]model.Value{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	result := ToProtoPart(part)

	assert.Equal(s.T(), inventoryV1.Category_CATEGORY_UNSPECIFIED, result.Category)
}

func (s *ConverterSuite) TestToProtoPart_AllMetadataTypes() {
	part := &model.Part{
		OrderUuid:    uuid.New(),
		Name:         "Test Part",
		Category:     model.ENGINE,
		Dimensions:   model.Dimensions{},
		Manufacturer: model.Manufacturer{},
		Tags:         []string{},
		Metadata: map[string]model.Value{
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := ToProtoPart(part)

	// Проверка всех типов метаданных
	assert.Equal(s.T(), "test string", result.Metadata["string_val"].GetStringValue())
	assert.Equal(s.T(), int64(42), result.Metadata["int_val"].GetInt64Value())
	assert.Equal(s.T(), 3.14, result.Metadata["float_val"].GetDoubleValue())
	assert.Equal(s.T(), true, result.Metadata["bool_val"].GetBoolValue())
}

func (s *ConverterSuite) TestToProtoPart_EmptyMetadata() {
	part := &model.Part{
		OrderUuid:    uuid.New(),
		Name:         "Test Part",
		Category:     model.ENGINE,
		Dimensions:   model.Dimensions{},
		Manufacturer: model.Manufacturer{},
		Tags:         []string{},
		Metadata:     map[string]model.Value{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	result := ToProtoPart(part)

	assert.Empty(s.T(), result.Metadata)
}

func (s *ConverterSuite) TestToProtoPart_ZeroValues() {
	part := &model.Part{
		OrderUuid:     uuid.New(),
		Name:          "",
		Description:   "",
		Price:         0,
		StockQuantity: 0,
		Category:      model.UNKNOWN,
		Dimensions: model.Dimensions{
			Length: 0,
			Width:  0,
			Height: 0,
			Weight: 0,
		},
		Manufacturer: model.Manufacturer{
			Name:    "",
			Country: "",
			Website: "",
		},
		Tags:      []string{},
		Metadata:  map[string]model.Value{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	result := ToProtoPart(part)

	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "", result.Name)
	assert.Equal(s.T(), "", result.Description)
	assert.Equal(s.T(), float64(0), result.Price)
	assert.Equal(s.T(), int64(0), result.StockQuantity)
	assert.Equal(s.T(), inventoryV1.Category_CATEGORY_UNSPECIFIED, result.Category)
	assert.Equal(s.T(), float64(0), result.Dimensions.Length)
	assert.Equal(s.T(), "", result.Manufacturer.Name)
	assert.Empty(s.T(), result.Tags)
	assert.Empty(s.T(), result.Metadata)
}
