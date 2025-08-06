package converter

import (
	"github.com/google/uuid"
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
	"github.com/stretchr/testify/assert"
)

func (s *ConverterSuite) TestToModelPart_CompleteFilter() {
	// Подготовка
	partUUID1 := uuid.New()
	partUUID2 := uuid.New()

	req := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			PartUuid: []string{
				partUUID1.String(),
				partUUID2.String(),
			},
			PartName: []string{
				"Rocket Engine",
				"Fuel Tank",
			},
			Category: []inventoryV1.Category{
				inventoryV1.Category_CATEGORY_ENGINE,
				inventoryV1.Category_CATEGORY_FUEL,
			},
			ManufacturerCountry: []string{
				"USA",
				"Russia",
			},
			Tags: []string{
				"engine",
				"fuel",
			},
		},
	}

	// Выполнение
	result := ToModelPart(req)

	// Проверка
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result.Uuids, 2)
	assert.Contains(s.T(), result.Uuids, partUUID1)
	assert.Contains(s.T(), result.Uuids, partUUID2)

	assert.Len(s.T(), result.Names, 2)
	assert.Contains(s.T(), result.Names, "Rocket Engine")
	assert.Contains(s.T(), result.Names, "Fuel Tank")

	assert.Len(s.T(), result.Categories, 2)
	assert.Contains(s.T(), result.Categories, model.ENGINE)
	assert.Contains(s.T(), result.Categories, model.FUEL)

	assert.Len(s.T(), result.ManufacturerCountries, 2)
	assert.Contains(s.T(), result.ManufacturerCountries, "USA")
	assert.Contains(s.T(), result.ManufacturerCountries, "Russia")

	assert.Len(s.T(), result.Tags, 2)
	assert.Contains(s.T(), result.Tags, "engine")
	assert.Contains(s.T(), result.Tags, "fuel")
}

func (s *ConverterSuite) TestToModelPart_AllCategories() {
	req := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Category: []inventoryV1.Category{
				inventoryV1.Category_CATEGORY_ENGINE,
				inventoryV1.Category_CATEGORY_FUEL,
				inventoryV1.Category_CATEGORY_PORTHOLE,
				inventoryV1.Category_CATEGORY_WING,
			},
		},
	}

	result := ToModelPart(req)

	assert.Len(s.T(), result.Categories, 4)
	assert.Contains(s.T(), result.Categories, model.ENGINE)
	assert.Contains(s.T(), result.Categories, model.FUEL)
	assert.Contains(s.T(), result.Categories, model.PORTHOLE)
	assert.Contains(s.T(), result.Categories, model.WING)
}

func (s *ConverterSuite) TestToModelPart_InvalidUUID() {
	req := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			PartUuid: []string{
				"invalid-uuid",
				uuid.New().String(), // валидный UUID
			},
		},
	}

	result := ToModelPart(req)

	// Должен содержать только валидный UUID
	assert.Len(s.T(), result.Uuids, 1)
}

func (s *ConverterSuite) TestToModelPart_EmptyFilter() {
	req := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{},
	}

	result := ToModelPart(req)

	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.Uuids)
	assert.Empty(s.T(), result.Names)
	assert.Empty(s.T(), result.Categories)
	assert.Empty(s.T(), result.ManufacturerCountries)
	assert.Empty(s.T(), result.Tags)
}

func (s *ConverterSuite) TestToModelPart_NilFilter() {
	req := &inventoryV1.ListPartsRequest{
		Filter: nil,
	}

	result := ToModelPart(req)

	assert.Nil(s.T(), result)
}

func (s *ConverterSuite) TestToModelPart_OnlyUUIDs() {
	partUUID := uuid.New()
	req := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			PartUuid: []string{partUUID.String()},
		},
	}

	result := ToModelPart(req)

	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result.Uuids, 1)
	assert.Equal(s.T(), partUUID, result.Uuids[0])
	assert.Empty(s.T(), result.Names)
	assert.Empty(s.T(), result.Categories)
	assert.Empty(s.T(), result.ManufacturerCountries)
	assert.Empty(s.T(), result.Tags)
}

func (s *ConverterSuite) TestToModelPart_OnlyNames() {
	req := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			PartName: []string{"Engine", "Fuel"},
		},
	}

	result := ToModelPart(req)

	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.Uuids)
	assert.Len(s.T(), result.Names, 2)
	assert.Contains(s.T(), result.Names, "Engine")
	assert.Contains(s.T(), result.Names, "Fuel")
	assert.Empty(s.T(), result.Categories)
	assert.Empty(s.T(), result.ManufacturerCountries)
	assert.Empty(s.T(), result.Tags)
}

func (s *ConverterSuite) TestToModelPart_OnlyCategories() {
	req := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Category: []inventoryV1.Category{
				inventoryV1.Category_CATEGORY_ENGINE,
			},
		},
	}

	result := ToModelPart(req)

	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.Uuids)
	assert.Empty(s.T(), result.Names)
	assert.Len(s.T(), result.Categories, 1)
	assert.Equal(s.T(), model.ENGINE, result.Categories[0])
	assert.Empty(s.T(), result.ManufacturerCountries)
	assert.Empty(s.T(), result.Tags)
}

func (s *ConverterSuite) TestProtoToModel_OnlyManufacturerCountries() {
	req := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			ManufacturerCountry: []string{"USA", "Germany"},
		},
	}

	result := ToModelPart(req)

	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.Uuids)
	assert.Empty(s.T(), result.Names)
	assert.Empty(s.T(), result.Categories)
	assert.Len(s.T(), result.ManufacturerCountries, 2)
	assert.Contains(s.T(), result.ManufacturerCountries, "USA")
	assert.Contains(s.T(), result.ManufacturerCountries, "Germany")
	assert.Empty(s.T(), result.Tags)
}

func (s *ConverterSuite) TestProtoToModel_OnlyTags() {
	req := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Tags: []string{"rocket", "space"},
		},
	}

	result := ToModelPart(req)

	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), result.Uuids)
	assert.Empty(s.T(), result.Names)
	assert.Empty(s.T(), result.Categories)
	assert.Empty(s.T(), result.ManufacturerCountries)
	assert.Len(s.T(), result.Tags, 2)
	assert.Contains(s.T(), result.Tags, "rocket")
	assert.Contains(s.T(), result.Tags, "space")
}

func (s *ConverterSuite) TestProtoToModel_MixedInvalidUUIDs() {
	validUUID := uuid.New()
	req := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			PartUuid: []string{
				"invalid-uuid-1",
				validUUID.String(),
				"another-invalid",
				"",
			},
		},
	}

	result := ToModelPart(req)

	// Должен содержать только валидный UUID
	assert.Len(s.T(), result.Uuids, 1)
	assert.Equal(s.T(), validUUID, result.Uuids[0])
}
