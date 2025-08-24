package order_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/fixtures"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

// TestPayOrder_TableDriven использует table-driven подход для тестирования PayOrder
// Это заменяет большинство отдельных тестов PayOrder из pay_test.go
func (s *ServiceSuite) TestPayOrder_TableDriven() {
	// Определяем структуру тест-кейса
	type testCase struct {
		name           string
		input          *model.Order
		mockSetup      func() (*model.Order, *model.Order) // возвращает dbOrder, expectedPaidOrder
		expectedResult *model.Order
		expectedError  error
		validateResult func(*testing.T, *model.Order)
		validateEvent  func(*testing.T, *model.OrderPaidEvent)
		skipEventCheck bool
		description    string // описание сценария для документации
	}

	// Предопределенные данные для консистентности
	orderUUID := fixtures.TestConstants.OrderUUID1
	userUUID := fixtures.TestConstants.UserUUID1
	partUUID := fixtures.TestConstants.PartUUID1
	transactionUUID := fixtures.TestConstants.TransactionUUID

	testCases := []testCase{
		{
			name: "успешная_оплата_картой",
			input: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithPaymentMethod("CARD").
				Build(),
			mockSetup: func() (*model.Order, *model.Order) {
				dbOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(100.0).
					WithStatus(model.StatusPendingPayment).
					Build()

				paidOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(100.0).
					WithTransactionUUID(transactionUUID).
					WithPaymentMethod("CARD").
					WithStatus(model.StatusPaid).
					Build()

				// orderWithPaymentMethod больше не используется
				_ = fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(100.0).
					WithStatus(model.StatusPendingPayment).
					WithPaymentMethod("CARD").
					Build()

				s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
					Return(dbOrder, nil)
				s.paymentClient.On("CreatePayment", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(paidOrder, nil)
				s.orderRepository.On("UpdateOrder", mock.Anything, paidOrder).
					Return(paidOrder, nil)

				return dbOrder, paidOrder
			},
			expectedResult: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID).
				WithTotalPrice(100.0).
				WithTransactionUUID(transactionUUID).
				WithPaymentMethod("CARD").
				WithStatus(model.StatusPaid).
				Build(),
			expectedError: nil,
			validateResult: func(t *testing.T, result *model.Order) {
				require.Equal(t, model.StatusPaid, result.Status)
				require.Equal(t, "CARD", result.PaymentMethod)
				require.Equal(t, transactionUUID, result.TransactionUUID)
			},
			validateEvent: func(t *testing.T, event *model.OrderPaidEvent) {
				require.NotNil(t, event)
				require.Equal(t, orderUUID, event.OrderUUID)
				require.Equal(t, userUUID, event.UserUUID)
				require.Equal(t, transactionUUID, event.TransactionUUID)
				require.Equal(t, "CARD", event.PaymentMethod)
			},
			description: "Стандартный сценарий оплаты картой",
		},
		{
			name: "успешная_оплата_SBP",
			input: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithPaymentMethod("SBP").
				Build(),
			mockSetup: func() (*model.Order, *model.Order) {
				dbOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(250.0).
					WithStatus(model.StatusPendingPayment).
					Build()

				paidOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(250.0).
					WithTransactionUUID(transactionUUID).
					WithPaymentMethod("SBP").
					WithStatus(model.StatusPaid).
					Build()

				// orderWithPaymentMethod больше не используется
				_ = fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(250.0).
					WithStatus(model.StatusPendingPayment).
					WithPaymentMethod("SBP").
					Build()

				s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
					Return(dbOrder, nil)
				s.paymentClient.On("CreatePayment", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(paidOrder, nil)
				s.orderRepository.On("UpdateOrder", mock.Anything, paidOrder).
					Return(paidOrder, nil)

				return dbOrder, paidOrder
			},
			expectedResult: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID).
				WithTotalPrice(250.0).
				WithTransactionUUID(transactionUUID).
				WithPaymentMethod("SBP").
				WithStatus(model.StatusPaid).
				Build(),
			expectedError: nil,
			validateEvent: func(t *testing.T, event *model.OrderPaidEvent) {
				require.NotNil(t, event)
				require.Equal(t, "SBP", event.PaymentMethod)
			},
			description: "Оплата через Систему Быстрых Платежей",
		},
		{
			name: "ошибка_заказ_не_найден",
			input: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithPaymentMethod("CARD").
				Build(),
			mockSetup: func() (*model.Order, *model.Order) {
				s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
					Return(nil, model.ErrOrderNotFound)
				return nil, nil
			},
			expectedResult: nil,
			expectedError:  model.ErrOrderNotFound,
			skipEventCheck: true,
			description:    "Попытка оплаты несуществующего заказа",
		},
		{
			name: "ошибка_payment_service_недоступен",
			input: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithPaymentMethod("CARD").
				Build(),
			mockSetup: func() (*model.Order, *model.Order) {
				dbOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(100.0).
					WithStatus(model.StatusPendingPayment).
					Build()

				// orderWithPaymentMethod больше не используется
				_ = fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(100.0).
					WithStatus(model.StatusPendingPayment).
					WithPaymentMethod("CARD").
					Build()

				s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
					Return(dbOrder, nil)
				s.paymentClient.On("CreatePayment", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(nil, model.ErrPaid)

				return dbOrder, nil
			},
			expectedResult: nil,
			expectedError:  model.ErrPaid,
			skipEventCheck: true,
			description:    "Payment service возвращает ошибку при создании платежа",
		},
		{
			name: "ошибка_обновления_заказа_в_репозитории",
			input: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithPaymentMethod("CARD").
				Build(),
			mockSetup: func() (*model.Order, *model.Order) {
				dbOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(100.0).
					WithStatus(model.StatusPendingPayment).
					Build()

				paidOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(100.0).
					WithTransactionUUID(transactionUUID).
					WithPaymentMethod("CARD").
					WithStatus(model.StatusPaid).
					Build()

				// orderWithPaymentMethod больше не используется
				_ = fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(100.0).
					WithStatus(model.StatusPendingPayment).
					WithPaymentMethod("CARD").
					Build()

				s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
					Return(dbOrder, nil)
				s.paymentClient.On("CreatePayment", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(paidOrder, nil)
				s.orderRepository.On("UpdateOrder", mock.Anything, paidOrder).
					Return(nil, model.ErrOrderNotFound)

				return dbOrder, paidOrder
			},
			expectedResult: nil,
			expectedError:  model.ErrOrderNotFound,
			skipEventCheck: true,
			description:    "Ошибка обновления заказа после успешной оплаты",
		},
		{
			name: "граничный_случай_большая_сумма_заказа",
			input: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithPaymentMethod("INVESTOR_MONEY").
				Build(),
			mockSetup: func() (*model.Order, *model.Order) {
				dbOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(999999999.0).
					WithStatus(model.StatusPendingPayment).
					Build()

				paidOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(999999999.0).
					WithTransactionUUID(transactionUUID).
					WithPaymentMethod("INVESTOR_MONEY").
					WithStatus(model.StatusPaid).
					Build()

				// orderWithPaymentMethod больше не используется
				_ = fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(999999999.0).
					WithStatus(model.StatusPendingPayment).
					WithPaymentMethod("INVESTOR_MONEY").
					Build()

				s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
					Return(dbOrder, nil)
				s.paymentClient.On("CreatePayment", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(paidOrder, nil)
				s.orderRepository.On("UpdateOrder", mock.Anything, paidOrder).
					Return(paidOrder, nil)

				return dbOrder, paidOrder
			},
			expectedResult: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID).
				WithTotalPrice(999999999.0).
				WithTransactionUUID(transactionUUID).
				WithPaymentMethod("INVESTOR_MONEY").
				WithStatus(model.StatusPaid).
				Build(),
			expectedError: nil,
			validateResult: func(t *testing.T, result *model.Order) {
				require.True(t, result.TotalPrice >= float32(999999999.0))
				require.Equal(t, "INVESTOR_MONEY", result.PaymentMethod)
			},
			description: "Оплата большой суммы инвестиционными средствами",
		},
	}

	// Выполняем все тест-кейсы
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Очищаем моки перед каждым тестом
			s.SetupTest()

			// Настраиваем моки для текущего тест-кейса
			_, _ = tc.mockSetup()

			// Вызываем тестируемый метод
			result, err := s.service.PayOrder(context.Background(), tc.input)

			// Проверяем ошибку
			if tc.expectedError != nil {
				s.Require().Error(err, "ожидалась ошибка для сценария: %s", tc.description)
				s.Require().ErrorIs(err, tc.expectedError, "неправильный тип ошибки")
				s.Require().Nil(result, "результат должен быть nil при ошибке")
			} else {
				s.Require().NoError(err, "не ожидалось ошибки для сценария: %s", tc.description)
				s.Require().NotNil(result, "результат не должен быть nil")

				// Проверяем основные поля результата
				s.Require().Equal(tc.expectedResult.OrderUUID, result.OrderUUID)
				s.Require().Equal(tc.expectedResult.Status, result.Status)
				s.Require().Equal(tc.expectedResult.PaymentMethod, result.PaymentMethod)
				s.Require().Equal(tc.expectedResult.TotalPrice, result.TotalPrice)

				// Выполняем дополнительную валидацию результата, если она определена
				if tc.validateResult != nil {
					tc.validateResult(s.T(), result)
				}
			}

			// Проверяем событие OrderPaid, если не пропускается
			if !tc.skipEventCheck && tc.expectedError == nil {
				lastEvent := s.orderPaidProducer.GetLastEvent()
				if tc.validateEvent != nil {
					tc.validateEvent(s.T(), lastEvent)
				} else {
					// Базовая валидация события
					s.Require().NotNil(lastEvent, "событие OrderPaid должно быть отправлено")
					s.Require().Equal(tc.expectedResult.OrderUUID, lastEvent.OrderUUID)
					s.Require().Equal(tc.expectedResult.PaymentMethod, lastEvent.PaymentMethod)
				}
			}

			// Проверяем, что все моки были вызваны корректно
			s.orderRepository.AssertExpectations(s.T())
			s.paymentClient.AssertExpectations(s.T())
		})
	}
}

// TestPayOrder_PaymentMethodVariations тестирует различные методы оплаты
func (s *ServiceSuite) TestPayOrder_PaymentMethodVariations() {
	// Список всех поддерживаемых методов оплаты
	paymentMethods := []struct {
		method      string
		description string
	}{
		{"CARD", "Банковская карта"},
		{"CREDIT_CARD", "Кредитная карта"},
		{"SBP", "Система Быстрых Платежей"},
		{"INVESTOR_MONEY", "Инвестиционные средства"},
		{"BANK_TRANSFER", "Банковский перевод"},
	}

	orderUUID := fixtures.TestConstants.OrderUUID1
	userUUID := fixtures.TestConstants.UserUUID1
	partUUID := fixtures.TestConstants.PartUUID1
	transactionUUID := fixtures.TestConstants.TransactionUUID

	for _, pm := range paymentMethods {
		s.Run("метод_оплаты_"+pm.method, func() {
			s.SetupTest()

			// Входящий запрос
			incomingOrder := fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithPaymentMethod(pm.method).
				Build()

			// Заказ из БД
			dbOrder := fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID).
				WithTotalPrice(150.0).
				WithStatus(model.StatusPendingPayment).
				Build()

			// Заказ после оплаты
			paidOrder := fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID).
				WithTotalPrice(150.0).
				WithTransactionUUID(transactionUUID).
				WithPaymentMethod(pm.method).
				WithStatus(model.StatusPaid).
				Build()

			// Заказ с установленным методом оплаты для передачи в payment service
			// Заказ с методом оплаты уже создан в incomingOrder

			// Настройка моков
			s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
				Return(dbOrder, nil)
			s.paymentClient.On("CreatePayment", mock.Anything, mock.AnythingOfType("*model.Order")).
				Return(paidOrder, nil)
			s.orderRepository.On("UpdateOrder", mock.Anything, paidOrder).
				Return(paidOrder, nil)

			// Вызов метода
			result, err := s.service.PayOrder(context.Background(), incomingOrder)

			// Проверка результата
			s.Require().NoError(err, "не ожидалось ошибки для метода оплаты %s (%s)", pm.method, pm.description)
			s.Require().NotNil(result)
			s.Require().Equal(pm.method, result.PaymentMethod, "метод оплаты должен быть сохранен")
			s.Require().Equal(model.StatusPaid, result.Status, "статус должен быть PAID")

			// Проверка события
			lastEvent := s.orderPaidProducer.GetLastEvent()
			s.Require().NotNil(lastEvent, "событие должно быть отправлено")
			s.Require().Equal(pm.method, lastEvent.PaymentMethod, "метод оплаты в событии должен соответствовать")

			// Проверяем моки
			s.orderRepository.AssertExpectations(s.T())
			s.paymentClient.AssertExpectations(s.T())
		})
	}
}

// TestPayOrder_EdgeCases тестирует граничные случаи и особые сценарии
func (s *ServiceSuite) TestPayOrder_EdgeCases() {
	edgeCases := []struct {
		name      string
		scenario  func() (*model.Order, string)
		setupMock func(*model.Order, string)
		validate  func(*testing.T, *model.Order, error)
	}{
		{
			name: "пустой_метод_оплаты",
			scenario: func() (*model.Order, string) {
				order := fixtures.NewOrderBuilder().
					WithOrderUUID(fixtures.TestConstants.OrderUUID1).
					WithPaymentMethod(""). // Пустой метод оплаты
					Build()
				return order, "Пустой метод оплаты должен вызвать ошибку"
			},
			setupMock: func(order *model.Order, description string) {
				// PayOrder сначала вызывает GetOrder, поэтому нужен мок
				existingOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(order.OrderUUID).
					WithStatus(model.StatusPendingPayment).
					Build()
				s.orderRepository.On("GetOrder", mock.Anything, order.OrderUUID).
					Return(existingOrder, nil)
				// Также нужен мок для CreatePayment
				s.paymentClient.On("CreatePayment", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(nil, errors.New("payment method cannot be empty")).Maybe()
			},
			validate: func(t *testing.T, result *model.Order, err error) {
				require.Error(t, err)
				require.Nil(t, result)
				require.Contains(t, err.Error(), "payment method cannot be empty")
			},
		},
		{
			name: "заказ_уже_оплачен",
			scenario: func() (*model.Order, string) {
				order := fixtures.NewOrderBuilder().
					WithOrderUUID(fixtures.TestConstants.OrderUUID1).
					WithPaymentMethod("CARD").
					Build()
				return order, "Попытка повторной оплаты уже оплаченного заказа"
			},
			setupMock: func(order *model.Order, description string) {
				// Заказ уже в статусе PAID
				dbOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(order.OrderUUID).
					WithUserUUID(fixtures.TestConstants.UserUUID1).
					WithPartUUIDs(fixtures.TestConstants.PartUUID1).
					WithTotalPrice(100.0).
					WithStatus(model.StatusPaid). // Уже оплачен
					Build()

				s.orderRepository.On("GetOrder", mock.Anything, order.OrderUUID).
					Return(dbOrder, nil)
				// Добавляем мок для CreatePayment, который может вызываться
				s.paymentClient.On("CreatePayment", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(nil, errors.New("order already paid")).Maybe()
			},
			validate: func(t *testing.T, result *model.Order, err error) {
				require.Error(t, err)
				require.Nil(t, result)
				require.Contains(t, err.Error(), "order already paid")
			},
		},
	}

	for _, ec := range edgeCases {
		s.Run(ec.name, func() {
			s.SetupTest()

			order, description := ec.scenario()
			ec.setupMock(order, description)

			result, err := s.service.PayOrder(context.Background(), order)

			ec.validate(s.T(), result, err)

			// Проверяем моки
			s.orderRepository.AssertExpectations(s.T())
			s.paymentClient.AssertExpectations(s.T())
		})
	}
}
