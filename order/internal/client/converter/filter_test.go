package converter

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

func (s *ConverterSuite) TestToProtoFilter_CompleteFilter() {
	// Подготовка
	partUUID1 := uuid.New()
	partUUID2 := uuid.New()
	filter := &model.Filter{
		PartUUIDs: []uuid.UUID{
			partUUID1,
			partUUID2,
		},
		PartNames: []string{
			"Rocket Engine",
			"Fuel Tank",
		},
		Categories: []model.Category{
			model.ENGINE,
			model.FUEL,
		},
		ManufacturerCountries: []string{
			"USA",
			"Russia",
		},
		Tags: []string{
			"engine",
			"fuel",
		},
	}

	// Выполнение
	result := ToProtoFilter(filter)

	// Проверка
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result.PartUuid, 2)
	assert.Contains(s.T(), result.PartUuid, partUUID1.String())
	assert.Contains(s.T(), result.PartUuid, partUUID2.String())

	assert.Len(s.T(), result.PartName, 2)
	assert.Contains(s.T(), result.PartName, "Rocket Engine")
	assert.Contains(s.T(), result.PartName, "Fuel Tank")

	assert.Len(s.T(), result.Category, 2)
	assert.Contains(s.T(), result.Category, inventoryV1.Category_CATEGORY_ENGINE)
	assert.Contains(s.T(), result.Category, inventoryV1.Category_CATEGORY_FUEL)

	assert.Len(s.T(), result.ManufacturerCountry, 2)
	assert.Contains(s.T(), result.ManufacturerCountry, "USA")
	assert.Contains(s.T(), result.ManufacturerCountry, "Russia")

	assert.Len(s.T(), result.Tags, 2)
	assert.Contains(s.T(), result.Tags, "engine")
	assert.Contains(s.T(), result.Tags, "fuel")
}

func (s *ConverterSuite) TestToProtoFilter_AllCategories() {
	// Подготовка
	filter := &model.Filter{
		Categories: []model.Category{
			model.ENGINE,
			model.FUEL,
			model.PORTHOLE,
			model.WING,
		},
	}

	// Выполнение
	result := ToProtoFilter(filter)

	// Проверка
	assert.Len(s.T(), result.Category, 4)
	assert.Contains(s.T(), result.Category, inventoryV1.Category_CATEGORY_ENGINE)
	assert.Contains(s.T(), result.Category, inventoryV1.Category_CATEGORY_FUEL)
	assert.Contains(s.T(), result.Category, inventoryV1.Category_CATEGORY_PORTHOLE)
	assert.Contains(s.T(), result.Category, inventoryV1.Category_CATEGORY_WING)
}

func (s *ConverterSuite) TestToProtoFilter_EmptyFilter() {
	// Подготовка
	filter := &model.Filter{}

	// Выполнение
	result := ToProtoFilter(filter)

	// Проверка
	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.PartUuid)
	assert.Empty(s.T(), result.PartName)
	assert.Empty(s.T(), result.Category)
	assert.Empty(s.T(), result.ManufacturerCountry)
	assert.Empty(s.T(), result.Tags)
}

func (s *ConverterSuite) TestToProtoFilter_OnlyUUIDs() {
	// Подготовка
	partUUID := uuid.New()
	filter := &model.Filter{
		PartUUIDs: []uuid.UUID{partUUID},
	}

	// Выполнение
	result := ToProtoFilter(filter)

	// Проверка
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result.PartUuid, 1)
	assert.Equal(s.T(), partUUID.String(), result.PartUuid[0])
	assert.Empty(s.T(), result.PartName)
	assert.Empty(s.T(), result.Category)
	assert.Empty(s.T(), result.ManufacturerCountry)
	assert.Empty(s.T(), result.Tags)
}

func (s *ConverterSuite) TestToProtoFilter_OnlyNames() {
	// Подготовка
	filter := &model.Filter{
		PartNames: []string{"Engine", "Fuel"},
	}

	// Выполнение
	result := ToProtoFilter(filter)

	// Проверка
	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.PartUuid)
	assert.Len(s.T(), result.PartName, 2)
	assert.Contains(s.T(), result.PartName, "Engine")
	assert.Contains(s.T(), result.PartName, "Fuel")
	assert.Empty(s.T(), result.Category)
	assert.Empty(s.T(), result.ManufacturerCountry)
	assert.Empty(s.T(), result.Tags)
}

func (s *ConverterSuite) TestToProtoFilter_OnlyCategories() {
	// Подготовка
	filter := &model.Filter{
		Categories: []model.Category{model.ENGINE},
	}

	// Выполнение
	result := ToProtoFilter(filter)

	// Проверка
	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.PartUuid)
	assert.Empty(s.T(), result.PartName)
	assert.Len(s.T(), result.Category, 1)
	assert.Equal(s.T(), inventoryV1.Category_CATEGORY_ENGINE, result.Category[0])
	assert.Empty(s.T(), result.ManufacturerCountry)
	assert.Empty(s.T(), result.Tags)
}

func (s *ConverterSuite) TestToProtoFilter_OnlyManufacturerCountries() {
	// Подготовка
	filter := &model.Filter{
		ManufacturerCountries: []string{"USA", "Germany"},
	}

	// Выполнение
	result := ToProtoFilter(filter)

	// Проверка
	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.PartUuid)
	assert.Empty(s.T(), result.PartName)
	assert.Empty(s.T(), result.Category)
	assert.Len(s.T(), result.ManufacturerCountry, 2)
	assert.Contains(s.T(), result.ManufacturerCountry, "USA")
	assert.Contains(s.T(), result.ManufacturerCountry, "Germany")
	assert.Empty(s.T(), result.Tags)
}

func (s *ConverterSuite) TestToProtoFilter_OnlyTags() {
	// Подготовка
	filter := &model.Filter{
		Tags: []string{"rocket", "space"},
	}

	// Выполнение
	result := ToProtoFilter(filter)

	// Проверка
	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.PartUuid)
	assert.Empty(s.T(), result.PartName)
	assert.Empty(s.T(), result.Category)
	assert.Empty(s.T(), result.ManufacturerCountry)
	assert.Len(s.T(), result.Tags, 2)
	assert.Contains(s.T(), result.Tags, "rocket")
	assert.Contains(s.T(), result.Tags, "space")
}
