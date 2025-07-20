package converter

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

func (s *ConverterSuite) TestPartToModel_Success() {
	// Подготовка
	partUUID := uuid.New()
	protoPart := &inventoryV1.Part{
		PartUuid:    partUUID.String(),
		Name:        "Rocket Engine",
		Description: "Powerful rocket engine",
		Price:       1500.75,
	}

	// Выполнение
	result, err := PartToModel(protoPart)

	// Проверка
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), partUUID, result.PartUUID)
	assert.Equal(s.T(), "Rocket Engine", result.Name)
	assert.Equal(s.T(), "Powerful rocket engine", result.Description)
	assert.Equal(s.T(), 1500.75, result.Price)
}

func (s *ConverterSuite) TestPartToModel_InvalidUUID() {
	// Подготовка
	protoPart := &inventoryV1.Part{
		PartUuid:    "invalid-uuid",
		Name:        "Rocket Engine",
		Description: "Powerful rocket engine",
		Price:       1500.75,
	}

	// Выполнение
	result, err := PartToModel(protoPart)

	// Проверка
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *ConverterSuite) TestPartToModel_EmptyUUID() {
	// Подготовка
	protoPart := &inventoryV1.Part{
		PartUuid:    "",
		Name:        "Rocket Engine",
		Description: "Powerful rocket engine",
		Price:       1500.75,
	}

	// Выполнение
	result, err := PartToModel(protoPart)

	// Проверка
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *ConverterSuite) TestPartToModel_ZeroValues() {
	// Подготовка
	partUUID := uuid.New()
	protoPart := &inventoryV1.Part{
		PartUuid:    partUUID.String(),
		Name:        "",
		Description: "",
		Price:       0,
	}

	// Выполнение
	result, err := PartToModel(protoPart)

	// Проверка
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), partUUID, result.PartUUID)
	assert.Equal(s.T(), "", result.Name)
	assert.Equal(s.T(), "", result.Description)
	assert.Equal(s.T(), float64(0), result.Price)
}

func (s *ConverterSuite) TestPartListToModel_Success() {
	// Подготовка
	partUUID1 := uuid.New()
	partUUID2 := uuid.New()
	protoParts := []*inventoryV1.Part{
		{
			PartUuid:    partUUID1.String(),
			Name:        "Rocket Engine",
			Description: "Powerful rocket engine",
			Price:       1500.75,
		},
		{
			PartUuid:    partUUID2.String(),
			Name:        "Fuel Tank",
			Description: "Large fuel tank",
			Price:       800.50,
		},
	}

	// Выполнение
	result, err := PartListToModel(protoParts)

	// Проверка
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), *result, 2)

	// Проверка первого элемента
	assert.Equal(s.T(), partUUID1, (*result)[0].PartUUID)
	assert.Equal(s.T(), "Rocket Engine", (*result)[0].Name)
	assert.Equal(s.T(), "Powerful rocket engine", (*result)[0].Description)
	assert.Equal(s.T(), 1500.75, (*result)[0].Price)

	// Проверка второго элемента
	assert.Equal(s.T(), partUUID2, (*result)[1].PartUUID)
	assert.Equal(s.T(), "Fuel Tank", (*result)[1].Name)
	assert.Equal(s.T(), "Large fuel tank", (*result)[1].Description)
	assert.Equal(s.T(), 800.50, (*result)[1].Price)
}

func (s *ConverterSuite) TestPartListToModel_EmptyList() {
	// Подготовка
	protoParts := []*inventoryV1.Part{}

	// Выполнение
	result, err := PartListToModel(protoParts)

	// Проверка
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), *result, 0)
}

func (s *ConverterSuite) TestPartListToModel_WithInvalidUUID() {
	// Подготовка
	partUUID := uuid.New()
	protoParts := []*inventoryV1.Part{
		{
			PartUuid:    partUUID.String(),
			Name:        "Rocket Engine",
			Description: "Powerful rocket engine",
			Price:       1500.75,
		},
		{
			PartUuid:    "invalid-uuid",
			Name:        "Fuel Tank",
			Description: "Large fuel tank",
			Price:       800.50,
		},
	}

	// Выполнение
	result, err := PartListToModel(protoParts)

	// Проверка
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}
