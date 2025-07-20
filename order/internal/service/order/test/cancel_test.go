package order_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	s.NoError(err)
	s.NotNil(result)
	s.Equal(model.StatusCancelled, result.Status)
	s.Equal(expectedOrder.OrderUUID, result.OrderUUID)

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
		Return(nil, status.Error(codes.NotFound, "order not found"))

	// Вызов метода
	result, err := s.service.CancelOrder(context.Background(), order)

	// Проверка результата
	s.Error(err)
	s.Nil(result)
	s.Equal(codes.NotFound, status.Code(err))
	s.Contains(status.Convert(err).Message(), "order not found")

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
	s.Error(err)
	s.Nil(result)
	s.Equal(codes.FailedPrecondition, status.Code(err))
	s.Contains(status.Convert(err).Message(), "order status is paid")

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
	s.Error(err)
	s.Nil(result)
	s.Equal(codes.FailedPrecondition, status.Code(err))
	s.Contains(status.Convert(err).Message(), "order status is cancelled")

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCancelOrder_UpdateError() {
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
		Status:     model.StatusPendingPayment,
	}

	// Настройка моков
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(existingOrder, nil)
	s.orderRepository.On("UpdateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
		Return(nil, status.Error(codes.Internal, "database error"))

	// Вызов метода
	result, err := s.service.CancelOrder(context.Background(), order)

	// Проверка результата
	s.Error(err)
	s.Nil(result)
	s.Equal(codes.Internal, status.Code(err))
	s.Contains(status.Convert(err).Message(), "failed to update order")

	s.orderRepository.AssertExpectations(s.T())
}
