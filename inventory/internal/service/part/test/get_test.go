package part_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	// Тестовые данные
	partUUID := uuid.New()
	expectedPart := &model.Part{
		OrderUuid:     partUUID,
		Name:          "Test Part",
		Description:   "Test Description",
		Price:         100.50,
		StockQuantity: 10,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	s.inventoryRepo.On("GetPart", s.ctx, partUUID).Return(expectedPart, nil)

	// Вызов метода
	result, err := s.service.GetPart(s.ctx, partUUID)

	// Проверка результата
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedPart, result)
	assert.Equal(s.T(), partUUID, result.OrderUuid)
	assert.Equal(s.T(), "Test Part", result.Name)
	s.inventoryRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetRepoError() {
	// Тестовые данные
	partUUID := uuid.New()
	expectedError := model.ErrPartNotFound

	s.inventoryRepo.On("GetPart", s.ctx, partUUID).Return(nil, expectedError)

	// 	Вызов метода
	result, err := s.service.GetPart(s.ctx, partUUID)

	// Проверка результата
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), expectedError, err)
	s.inventoryRepo.AssertExpectations(s.T())
}
