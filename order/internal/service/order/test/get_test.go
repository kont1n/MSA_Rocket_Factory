package order_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s *ServiceSuite) TestGetOrder_Success() {
	// Тестовые данные
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUID := uuid.New()

	expectedOrder := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  []uuid.UUID{partUUID},
		TotalPrice: 100.0,
		Status:     model.StatusPendingPayment,
	}

	// Настройка моков
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(expectedOrder, nil)

	// Вызов метода
	result, err := s.service.GetOrder(context.Background(), orderUUID)

	// Проверка результата
	s.NoError(err)
	s.NotNil(result)
	s.Equal(expectedOrder.OrderUUID, result.OrderUUID)
	s.Equal(expectedOrder.UserUUID, result.UserUUID)
	s.Equal(expectedOrder.TotalPrice, result.TotalPrice)
	s.Equal(expectedOrder.Status, result.Status)

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetOrder_NotFound() {
	// Тестовые данные
	orderUUID := uuid.New()

	// Настройка моков - симулируем отсутствие заказа
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(nil, status.Error(codes.NotFound, "order not found"))

	// Вызов метода
	result, err := s.service.GetOrder(context.Background(), orderUUID)

	// Проверка результата
	s.Error(err)
	s.Nil(result)
	s.Equal(codes.FailedPrecondition, status.Code(err))
	s.Contains(status.Convert(err).Message(), "order not found")

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetOrder_RepositoryError() {
	// Тестовые данные
	orderUUID := uuid.New()

	// Настройка моков - симулируем ошибку репозитория
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(nil, status.Error(codes.Internal, "database error"))

	// Вызов метода
	result, err := s.service.GetOrder(context.Background(), orderUUID)

	// Проверка результата
	s.Error(err)
	s.Nil(result)
	s.Equal(codes.FailedPrecondition, status.Code(err))
	s.Contains(status.Convert(err).Message(), "order not found")

	s.orderRepository.AssertExpectations(s.T())
}
