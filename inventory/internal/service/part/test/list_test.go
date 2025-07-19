package part_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
)

func (s *ServiceSuite) TestListSuccess() {
	// Тестовые данные
	filter := &model.Filter{
		Names: []string{"Test Part"},
	}

	expectedParts := []model.Part{
		{
			OrderUuid:     uuid.New(),
			Name:          "Test Part 1",
			Description:   "Test Description 1",
			Price:         100.50,
			StockQuantity: 10,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			OrderUuid:     uuid.New(),
			Name:          "Test Part 2",
			Description:   "Test Description 2",
			Price:         200.75,
			StockQuantity: 5,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	// Настраиваем мок
	s.inventoryRepo.On("ListParts", s.ctx, filter).Return(&expectedParts, nil)

	// Вызов метода
	result, err := s.service.ListParts(s.ctx, filter)

	// Проверка результата
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), *result, 2)
	assert.Equal(s.T(), "Test Part 1", (*result)[0].Name)
	assert.Equal(s.T(), "Test Part 2", (*result)[1].Name)

	// Проверяем, что мок был вызван
	s.inventoryRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListEmptyResult() {
	// Тестовые данные
	filter := &model.Filter{
		Names: []string{"NonExistentPart"},
	}

	expectedParts := []model.Part{}

	s.inventoryRepo.On("ListParts", s.ctx, filter).Return(&expectedParts, nil)

	// Вызов метода
	result, err := s.service.ListParts(s.ctx, filter)

	// Проверка результата
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Empty(s.T(), *result)
	s.inventoryRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListRepoError() {
	// Тестовые данные
	filter := &model.Filter{
		Names: []string{"Test Part"},
	}

	expectedError := assert.AnError

	s.inventoryRepo.On("ListParts", s.ctx, filter).Return(nil, expectedError)

	// Вызов метода
	result, err := s.service.ListParts(s.ctx, filter)

	// Проверка результата
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), expectedError, err)
	s.inventoryRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListWithNilFilter() {
	// Тестовые данные
	expectedParts := []model.Part{
		{
			OrderUuid:     uuid.New(),
			Name:          "All Parts",
			Description:   "All parts without filter",
			Price:         150.00,
			StockQuantity: 15,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	s.inventoryRepo.On("ListParts", s.ctx, (*model.Filter)(nil)).Return(&expectedParts, nil)

	// Вызов метода
	result, err := s.service.ListParts(s.ctx, nil)

	// Проверка результата
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), *result, 1)
	s.inventoryRepo.AssertExpectations(s.T())
}
