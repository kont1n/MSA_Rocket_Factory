package order_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s *ServiceSuite) TestCreateOrder_Success() {
	// Тестовые данные
	userUUID := uuid.New()
	partUUID1 := uuid.New()
	partUUID2 := uuid.New()

	order := &model.Order{
		UserUUID:  userUUID,
		PartUUIDs: []uuid.UUID{partUUID1, partUUID2},
	}

	parts := []model.Part{
		{PartUUID: partUUID1, Price: 100.0},
		{PartUUID: partUUID2, Price: 200.0},
	}

	expectedOrder := &model.Order{
		UserUUID:   userUUID,
		PartUUIDs:  []uuid.UUID{partUUID1, partUUID2},
		TotalPrice: 300.0,
		Status:     model.StatusPendingPayment,
	}

	// Настройка моков
	s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
		Return(&parts, nil)
	s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
		Return(expectedOrder, nil)

	// Вызов метода
	result, err := s.service.CreateOrder(context.Background(), order)

	// Проверка результата
	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Require().Equal(expectedOrder.TotalPrice, result.TotalPrice)
	s.Require().Equal(expectedOrder.Status, result.Status)

	s.inventoryClient.AssertExpectations(s.T())
	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCreateOrder_EmptyParts() {
	// Тестовые данные
	order := &model.Order{
		UserUUID:  uuid.New(),
		PartUUIDs: []uuid.UUID{},
	}

	// Вызов метода
	result, err := s.service.CreateOrder(context.Background(), order)

	// Проверка результата
	s.Require().Empty(result)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrPartsSpecified)
}

func (s *ServiceSuite) TestCreateOrder_PartsNotFound() {
	// Тестовые данные
	userUUID := uuid.New()
	partUUID1 := uuid.New()
	partUUID2 := uuid.New()

	order := &model.Order{
		UserUUID:  userUUID,
		PartUUIDs: []uuid.UUID{partUUID1, partUUID2},
	}

	parts := []model.Part{
		{PartUUID: partUUID1, Price: 100.0},
	}

	// Настройка моков - возвращаем только одну деталь вместо двух
	s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
		Return(&parts, nil)

	// Вызов метода
	result, err := s.service.CreateOrder(context.Background(), order)

	// Проверка результата
	s.Require().Empty(result)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrPartsListNotFound)

	s.inventoryClient.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCreateOrder_InventoryError() {
	// Тестовые данные
	order := &model.Order{
		UserUUID:  uuid.New(),
		PartUUIDs: []uuid.UUID{uuid.New()},
	}

	// Настройка моков - симулируем ошибку inventory
	s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
		Return(nil, status.Error(codes.Internal, "inventory service error"))

	// Вызов метода
	result, err := s.service.CreateOrder(context.Background(), order)

	// Проверка результата
	s.Require().Empty(result)
	s.Require().Error(err)
	s.Require().ErrorIs(err, status.Error(codes.Internal, "inventory service error"))

	s.inventoryClient.AssertExpectations(s.T())
}
