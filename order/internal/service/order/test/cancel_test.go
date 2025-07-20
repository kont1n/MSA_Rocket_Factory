package order_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s *ServiceSuite) TestCancelOrder_Success() {
	// Тестовые данные
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUID := uuid.New()

	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  []uuid.UUID{partUUID},
		TotalPrice: 100.0,
		Status:     model.StatusPendingPayment,
	}

	expectedOrder := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  []uuid.UUID{partUUID},
		TotalPrice: 100.0,
		Status:     model.StatusCancelled,
	}

	// Настройка моков
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(order, nil)
	s.orderRepository.On("UpdateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
		Return(expectedOrder, nil)

	// Вызов метода
	result, err := s.service.CancelOrder(context.Background(), order)

	// Проверка результата
	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Require().Equal(model.StatusCancelled, result.Status)
	s.Require().Equal(expectedOrder.OrderUUID, result.OrderUUID)

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCancelOrder_OrderNotFound() {
	// Тестовые данные
	orderUUID := uuid.New()
	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   uuid.New(),
		PartUUIDs:  []uuid.UUID{uuid.New()},
		TotalPrice: 100.0,
		Status:     model.StatusPendingPayment,
	}

	// Настройка моков - симулируем отсутствие заказа
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(nil, model.ErrOrderNotFound)

	// Вызов метода
	result, err := s.service.CancelOrder(context.Background(), order)

	// Проверка результата
	s.Require().Empty(result)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCancelOrder_AlreadyPaid() {
	// Тестовые данные
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUID := uuid.New()

	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  []uuid.UUID{partUUID},
		TotalPrice: 100.0,
		Status:     model.StatusPendingPayment,
	}

	existingOrder := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  []uuid.UUID{partUUID},
		TotalPrice: 100.0,
		Status:     model.StatusPaid,
	}

	// Настройка моков
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(existingOrder, nil)

	// Вызов метода
	result, err := s.service.CancelOrder(context.Background(), order)

	// Проверка результата
	s.Require().Empty(result)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrPaid)

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCancelOrder_AlreadyCancelled() {
	// Тестовые данные
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUID := uuid.New()

	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  []uuid.UUID{partUUID},
		TotalPrice: 100.0,
		Status:     model.StatusPendingPayment,
	}

	existingOrder := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  []uuid.UUID{partUUID},
		TotalPrice: 100.0,
		Status:     model.StatusCancelled,
	}

	// Настройка моков
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(existingOrder, nil)

	// Вызов метода
	result, err := s.service.CancelOrder(context.Background(), order)

	// Проверка результата
	s.Require().Empty(result)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrCancelled)

	s.orderRepository.AssertExpectations(s.T())
}
