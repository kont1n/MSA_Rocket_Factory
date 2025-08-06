package order_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/stretchr/testify/mock"
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
	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Require().Equal(expectedOrder.OrderUUID, result.OrderUUID)
	s.Require().Equal(expectedOrder.UserUUID, result.UserUUID)
	s.Require().Equal(expectedOrder.TotalPrice, result.TotalPrice)
	s.Require().Equal(expectedOrder.Status, result.Status)

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetOrder_NotFound() {
	// Тестовые данные
	orderUUID := uuid.New()

	// Настройка моков - симулируем отсутствие заказа
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(nil, model.ErrOrderNotFound)

	// Вызов метода
	result, err := s.service.GetOrder(context.Background(), orderUUID)

	// Проверка результата
	s.Require().Empty(result)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)

	s.orderRepository.AssertExpectations(s.T())
}
