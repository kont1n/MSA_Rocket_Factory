package order_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
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
	result, err := s.service.PayOrder(s.ctx, order)

	// Проверка результата
	s.NoError(err)
	s.NotNil(result)
	s.Equal(model.StatusPaid, result.Status)
	s.Equal(transactionUUID, result.TransactionUUID)
	s.Equal("CARD", result.PaymentMethod)

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
		Return(nil, status.Error(codes.NotFound, "order not found"))

	// Вызов метода
	result, err := s.service.PayOrder(s.ctx, order)

	// Проверка результата
	s.Error(err)
	s.Nil(result)
	s.Equal(codes.NotFound, status.Code(err))
	s.Contains(status.Convert(err).Message(), "order not found")

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestPayOrder_PaymentError() {
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
	s.paymentClient.On("CreatePayment", mock.Anything, existingOrder).
		Return(nil, status.Error(codes.FailedPrecondition, "payment failed"))

	// Вызов метода
	result, err := s.service.PayOrder(s.ctx, order)

	// Проверка результата
	s.Error(err)
	s.Nil(result)
	s.Equal(codes.FailedPrecondition, status.Code(err))
	s.Contains(status.Convert(err).Message(), "failed to pay for the order")

	s.orderRepository.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestPayOrder_UpdateError() {
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

	existingOrder := &model.Order{
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
		Return(existingOrder, nil)
	s.paymentClient.On("CreatePayment", mock.Anything, existingOrder).
		Return(paidOrder, nil)
	s.orderRepository.On("UpdateOrder", mock.Anything, paidOrder).
		Return(nil, status.Error(codes.Internal, "repository error"))

	// Вызов метода
	result, err := s.service.PayOrder(s.ctx, order)

	// Проверка результата
	s.Error(err)
	s.Nil(result)
	s.Equal(codes.Internal, status.Code(err))
	s.Contains(status.Convert(err).Message(), "failed to update order")

	s.orderRepository.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
}
