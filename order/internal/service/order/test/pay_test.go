package order_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestPayOrder_Success() {
	// Тестовые данные
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUID := uuid.New()
	transactionUUID := uuid.New()

	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  []uuid.UUID{partUUID},
		TotalPrice: 100.0,
		Status:     model.StatusPendingPayment,
	}

	paidOrder := &model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       []uuid.UUID{partUUID},
		TotalPrice:      100.0,
		TransactionUUID: transactionUUID,
		PaymentMethod:   "CARD",
		Status:          model.StatusPaid,
	}

	// Настройка моков
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(order, nil)
	s.paymentClient.On("CreatePayment", mock.Anything, order).
		Return(paidOrder, nil)
	s.orderRepository.On("UpdateOrder", mock.Anything, paidOrder).
		Return(paidOrder, nil)

	// Вызов метода
	result, err := s.service.PayOrder(context.Background(), order)

	// Проверка результата
	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Require().Equal(model.StatusPaid, result.Status)
	s.Require().Equal(transactionUUID, result.TransactionUUID)
	s.Require().Equal("CARD", result.PaymentMethod)

	s.orderRepository.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestPayOrder_OrderNotFound() {
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
	result, err := s.service.PayOrder(context.Background(), order)

	// Проверка результата
	s.Require().Empty(result)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)

	s.orderRepository.AssertExpectations(s.T())
}

