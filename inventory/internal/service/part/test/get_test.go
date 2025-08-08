package part_test

import (
	"context"
	"errors"
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

	s.inventoryRepo.On("GetPart", context.Background(), partUUID).Return(expectedPart, nil)

	// Вызов метода
	result, err := s.service.GetPart(context.Background(), partUUID)

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

	s.inventoryRepo.On("GetPart", context.Background(), partUUID).Return(nil, expectedError)

	// 	Вызов метода
	result, err := s.service.GetPart(context.Background(), partUUID)

	// Проверка результата
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.True(s.T(), errors.Is(err, expectedError))
	s.inventoryRepo.AssertExpectations(s.T())
}
