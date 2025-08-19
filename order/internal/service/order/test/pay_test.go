package order_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s *ServiceSuite) TestPayOrder_Success() {
	// Тестовые данные
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUID := uuid.New()
	transactionUUID := uuid.New()

	// Входящий запрос с методом оплаты
	incomingOrder := &model.Order{
		OrderUUID:     orderUUID,
		PaymentMethod: "CARD",
	}

	// Заказ из БД (без метода оплаты)
	dbOrder := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  []uuid.UUID{partUUID},
		TotalPrice: 100.0,
		Status:     model.StatusPendingPayment,
	}

	// Заказ после оплаты
	paidOrder := &model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUUIDs:       []uuid.UUID{partUUID},
		TotalPrice:      100.0,
		TransactionUUID: transactionUUID,
		PaymentMethod:   "CARD", // Должен быть установлен из входящего запроса
		Status:          model.StatusPaid,
	}

	// Настройка моков
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(dbOrder, nil)

	// Создаем заказ с установленным PaymentMethod для передачи в CreatePayment
	orderWithPaymentMethod := &model.Order{
		OrderUUID:     orderUUID,
		UserUUID:      userUUID,
		PartUUIDs:     []uuid.UUID{partUUID},
		TotalPrice:    100.0,
		Status:        model.StatusPendingPayment,
		PaymentMethod: "CARD", // Установлен из входящего запроса
	}

	s.paymentClient.On("CreatePayment", mock.Anything, orderWithPaymentMethod).
		Return(paidOrder, nil)
	s.orderRepository.On("UpdateOrder", mock.Anything, paidOrder).
		Return(paidOrder, nil)

	// Вызов метода
	result, err := s.service.PayOrder(context.Background(), incomingOrder)

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
	incomingOrder := &model.Order{
		OrderUUID:     orderUUID,
		PaymentMethod: "CARD",
	}

	// Настройка моков - симулируем отсутствие заказа
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(nil, model.ErrOrderNotFound)

	// Вызов метода
	result, err := s.service.PayOrder(context.Background(), incomingOrder)

	// Проверка результата
	s.Require().Empty(result)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestPayOrder_PaymentError() {
	// Тестовые данные
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUID := uuid.New()

	incomingOrder := &model.Order{
		OrderUUID:     orderUUID,
		PaymentMethod: "CARD",
	}

	dbOrder := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUUIDs:  []uuid.UUID{partUUID},
		TotalPrice: 100.0,
		Status:     model.StatusPendingPayment,
	}

	// Настройка моков - успешное получение заказа, но ошибка при создании платежа
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(dbOrder, nil)

	orderWithPaymentMethod := &model.Order{
		OrderUUID:     orderUUID,
		UserUUID:      userUUID,
		PartUUIDs:     []uuid.UUID{partUUID},
		TotalPrice:    100.0,
		Status:        model.StatusPendingPayment,
		PaymentMethod: "CARD",
	}

	s.paymentClient.On("CreatePayment", mock.Anything, orderWithPaymentMethod).
		Return(nil, model.ErrPaid)

	// Вызов метода
	result, err := s.service.PayOrder(context.Background(), incomingOrder)

	// Проверка результата
	s.Require().Empty(result)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrPaid)

	s.orderRepository.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestPayOrder_UpdateError() {
	// Тестовые данные
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUID := uuid.New()
	transactionUUID := uuid.New()

	incomingOrder := &model.Order{
		OrderUUID:     orderUUID,
		PaymentMethod: "CARD",
	}

	dbOrder := &model.Order{
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

	// Настройка моков - успешное получение заказа и создание платежа, но ошибка при обновлении
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(dbOrder, nil)

	orderWithPaymentMethod := &model.Order{
		OrderUUID:     orderUUID,
		UserUUID:      userUUID,
		PartUUIDs:     []uuid.UUID{partUUID},
		TotalPrice:    100.0,
		Status:        model.StatusPendingPayment,
		PaymentMethod: "CARD",
	}

	s.paymentClient.On("CreatePayment", mock.Anything, orderWithPaymentMethod).
		Return(paidOrder, nil)
	s.orderRepository.On("UpdateOrder", mock.Anything, paidOrder).
		Return(nil, model.ErrOrderNotFound)

	// Вызов метода
	result, err := s.service.PayOrder(context.Background(), incomingOrder)

	// Проверка результата
	s.Require().Empty(result)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)

	s.orderRepository.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestPayOrder_PaymentMethodPropagation() {
	// Тестовые данные
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUID := uuid.New()
	transactionUUID := uuid.New()

	// Входящий запрос с различными методами оплаты
	testCases := []struct {
		name          string
		paymentMethod string
	}{
		{"Credit Card", "CREDIT_CARD"},
		{"SBP", "SBP"},
		{"Investor Money", "INVESTOR_MONEY"},
		{"Card", "CARD"},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			incomingOrder := &model.Order{
				OrderUUID:     orderUUID,
				PaymentMethod: tc.paymentMethod,
			}

			dbOrder := &model.Order{
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
				PaymentMethod:   tc.paymentMethod,
				Status:          model.StatusPaid,
			}

			// Настройка моков
			s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
				Return(dbOrder, nil)

			orderWithPaymentMethod := &model.Order{
				OrderUUID:     orderUUID,
				UserUUID:      userUUID,
				PartUUIDs:     []uuid.UUID{partUUID},
				TotalPrice:    100.0,
				Status:        model.StatusPendingPayment,
				PaymentMethod: tc.paymentMethod,
			}

			s.paymentClient.On("CreatePayment", mock.Anything, orderWithPaymentMethod).
				Return(paidOrder, nil)
			s.orderRepository.On("UpdateOrder", mock.Anything, paidOrder).
				Return(paidOrder, nil)

			// Вызов метода
			result, err := s.service.PayOrder(context.Background(), incomingOrder)

			// Проверка результата
			s.Require().NoError(err)
			s.Require().NotNil(result)
			s.Require().Equal(tc.paymentMethod, result.PaymentMethod)
			s.Require().Equal(model.StatusPaid, result.Status)

			// Очищаем моки для следующего теста
			s.orderRepository.ExpectedCalls = nil
			s.paymentClient.ExpectedCalls = nil
		})
	}
}

func (s *ServiceSuite) TestPayOrder_EventPaymentMethodPropagation() {
	// Тестовые данные
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUID := uuid.New()
	transactionUUID := uuid.New()

	incomingOrder := &model.Order{
		OrderUUID:     orderUUID,
		PaymentMethod: "SBP",
	}

	dbOrder := &model.Order{
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
		PaymentMethod:   "SBP",
		Status:          model.StatusPaid,
	}

	// Настройка моков
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(dbOrder, nil)

	orderWithPaymentMethod := &model.Order{
		OrderUUID:     orderUUID,
		UserUUID:      userUUID,
		PartUUIDs:     []uuid.UUID{partUUID},
		TotalPrice:    100.0,
		Status:        model.StatusPendingPayment,
		PaymentMethod: "SBP",
	}

	s.paymentClient.On("CreatePayment", mock.Anything, orderWithPaymentMethod).
		Return(paidOrder, nil)
	s.orderRepository.On("UpdateOrder", mock.Anything, paidOrder).
		Return(paidOrder, nil)

	// Вызов метода
	result, err := s.service.PayOrder(context.Background(), incomingOrder)

	// Проверка результата
	s.Require().NoError(err)
	s.Require().NotNil(result)

	// Проверяем что событие было отправлено с правильным PaymentMethod
	lastEvent := s.orderPaidProducer.GetLastEvent()
	s.Require().NotNil(lastEvent)
	s.Require().Equal("SBP", lastEvent.PaymentMethod)
	s.Require().Equal(orderUUID, lastEvent.OrderUUID)
	s.Require().Equal(userUUID, lastEvent.UserUUID)
	s.Require().Equal(transactionUUID, lastEvent.TransactionUUID)

	s.orderRepository.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
}
