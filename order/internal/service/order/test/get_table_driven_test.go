package order_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/fixtures"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

// TestGetOrder_TableDriven использует table-driven подход для тестирования GetOrder
// Это заменяет отдельные тесты GetOrder из get_test.go
func (s *ServiceSuite) TestGetOrder_TableDriven() {
	type testCase struct {
		name           string
		orderUUID      uuid.UUID
		mockSetup      func(uuid.UUID)
		expectedResult *model.Order
		expectedError  error
		validateResult func(*testing.T, *model.Order)
		description    string
	}

	// Предопределенные данные
	validOrderUUID := fixtures.TestConstants.OrderUUID1
	invalidOrderUUID := fixtures.TestConstants.OrderUUID2
	userUUID := fixtures.TestConstants.UserUUID1
	partUUID := fixtures.TestConstants.PartUUID1

	testCases := []testCase{
		{
			name:      "успешное_получение_заказа",
			orderUUID: validOrderUUID,
			mockSetup: func(orderUUID uuid.UUID) {
				expectedOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(150.0).
					WithStatus(model.StatusPendingPayment).
					Build()

				s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
					Return(expectedOrder, nil)
			},
			expectedResult: fixtures.NewOrderBuilder().
				WithOrderUUID(validOrderUUID).
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID).
				WithTotalPrice(150.0).
				WithStatus(model.StatusPendingPayment).
				Build(),
			expectedError: nil,
			validateResult: func(t *testing.T, result *model.Order) {
				require.Equal(t, validOrderUUID, result.OrderUUID)
				require.Equal(t, userUUID, result.UserUUID)
				require.Equal(t, float32(150.0), result.TotalPrice)
				require.Equal(t, model.StatusPendingPayment, result.Status)
				require.Len(t, result.PartUUIDs, 1)
			},
			description: "Стандартный сценарий получения существующего заказа",
		},
		{
			name:      "успешное_получение_оплаченного_заказа",
			orderUUID: validOrderUUID,
			mockSetup: func(orderUUID uuid.UUID) {
				expectedOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(300.0).
					WithTransactionUUID(fixtures.TestConstants.TransactionUUID).
					WithStatus(model.StatusPaid).
					WithPaymentMethod("CARD").
					Build()

				s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
					Return(expectedOrder, nil)
			},
			expectedResult: fixtures.NewOrderBuilder().
				WithOrderUUID(validOrderUUID).
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID).
				WithTotalPrice(300.0).
				WithTransactionUUID(fixtures.TestConstants.TransactionUUID).
				WithStatus(model.StatusPaid).
				WithPaymentMethod("CARD").
				Build(),
			expectedError: nil,
			validateResult: func(t *testing.T, result *model.Order) {
				require.Equal(t, model.StatusPaid, result.Status)
				require.Equal(t, "CARD", result.PaymentMethod)
				require.NotEqual(t, uuid.Nil, result.TransactionUUID)
			},
			description: "Получение заказа в статусе оплачен",
		},
		{
			name:      "заказ_не_найден",
			orderUUID: invalidOrderUUID,
			mockSetup: func(orderUUID uuid.UUID) {
				s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
					Return(nil, model.ErrOrderNotFound)
			},
			expectedResult: nil,
			expectedError:  model.ErrOrderNotFound,
			description:    "Попытка получения несуществующего заказа",
		},
		{
			name:      "граничный_случай_заказ_с_большой_суммой",
			orderUUID: validOrderUUID,
			mockSetup: func(orderUUID uuid.UUID) {
				expectedOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(999999.99).
					WithStatus(model.StatusPendingPayment).
					Build()

				s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
					Return(expectedOrder, nil)
			},
			expectedResult: fixtures.NewOrderBuilder().
				WithOrderUUID(validOrderUUID).
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID).
				WithTotalPrice(999999.99).
				WithStatus(model.StatusPendingPayment).
				Build(),
			expectedError: nil,
			validateResult: func(t *testing.T, result *model.Order) {
				require.True(t, result.TotalPrice > float32(999999.0))
			},
			description: "Получение заказа с максимальной суммой",
		},
		{
			name:      "заказ_с_множественными_деталями",
			orderUUID: validOrderUUID,
			mockSetup: func(orderUUID uuid.UUID) {
				expectedOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(
						fixtures.TestConstants.PartUUID1,
						fixtures.TestConstants.PartUUID2,
						fixtures.TestConstants.PartUUID3,
					).
					WithTotalPrice(450.0).
					WithStatus(model.StatusPaid).
					Build()

				s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
					Return(expectedOrder, nil)
			},
			expectedResult: fixtures.NewOrderBuilder().
				WithOrderUUID(validOrderUUID).
				WithUserUUID(userUUID).
				WithPartUUIDs(
					fixtures.TestConstants.PartUUID1,
					fixtures.TestConstants.PartUUID2,
					fixtures.TestConstants.PartUUID3,
				).
				WithTotalPrice(450.0).
				WithStatus(model.StatusPaid).
				Build(),
			expectedError: nil,
			validateResult: func(t *testing.T, result *model.Order) {
				require.Len(t, result.PartUUIDs, 3)
				require.Equal(t, model.StatusPaid, result.Status)
				require.Equal(t, float32(450.0), result.TotalPrice)
			},
			description: "Получение заказа с несколькими деталями в статусе оплачен",
		},
	}

	// Выполняем все тест-кейсы
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Убираем t.Parallel() для избежания конфликтов
			s.SetupTest()

			// Настраиваем моки для текущего тест-кейса
			tc.mockSetup(tc.orderUUID)

			// Вызываем тестируемый метод
			result, err := s.service.GetOrder(context.Background(), tc.orderUUID)

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
				s.Require().Equal(tc.expectedResult.UserUUID, result.UserUUID)
				s.Require().Equal(tc.expectedResult.TotalPrice, result.TotalPrice)
				s.Require().Equal(tc.expectedResult.Status, result.Status)
				s.Require().Equal(len(tc.expectedResult.PartUUIDs), len(result.PartUUIDs))

				// Выполняем дополнительную валидацию результата, если она определена
				if tc.validateResult != nil {
					tc.validateResult(s.T(), result)
				}
			}

			// Проверяем, что все моки были вызваны корректно
			s.orderRepository.AssertExpectations(s.T())
		})
	}
}

// TestGetOrder_ConcurrentAccess тестирует конкурентный доступ к GetOrder
func (s *ServiceSuite) TestGetOrder_ConcurrentAccess() {
	s.T().Parallel()
	s.SetupTest() // Инициализируем моки для теста

	const numGoroutines = 20
	orderUUID := fixtures.TestConstants.OrderUUID1

	// Подготавливаем данные
	expectedOrder := fixtures.NewOrderBuilder().
		WithOrderUUID(orderUUID).
		WithUserUUID(fixtures.TestConstants.UserUUID1).
		WithPartUUIDs(fixtures.TestConstants.PartUUID1).
		WithTotalPrice(200.0).
		WithStatus(model.StatusPaid).
		Build()

	// Настраиваем моки для всех параллельных запросов
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(expectedOrder, nil).
		Times(numGoroutines)

	results := make(chan *model.Order, numGoroutines)
	errors := make(chan error, numGoroutines)

	// Запускаем горутины
	for i := 0; i < numGoroutines; i++ {
		go func() {
			result, err := s.service.GetOrder(context.Background(), orderUUID)
			results <- result
			errors <- err
		}()
	}

	// Собираем результаты
	for i := 0; i < numGoroutines; i++ {
		result := <-results
		err := <-errors

		s.Require().NoError(err, "не ожидалось ошибки в параллельном запросе")
		s.Require().NotNil(result, "результат не должен быть nil")
		s.Require().Equal(expectedOrder.OrderUUID, result.OrderUUID)
		s.Require().Equal(expectedOrder.TotalPrice, result.TotalPrice)
		s.Require().Equal(expectedOrder.Status, result.Status)
	}

	// Проверяем, что все вызовы моков произошли
	s.orderRepository.AssertExpectations(s.T())
}

// TestGetOrder_Performance простой тест производительности
func (s *ServiceSuite) TestGetOrder_Performance() {
	if testing.Short() {
		s.T().Skip("пропускаем тест производительности в short режиме")
	}

	orderUUID := fixtures.TestConstants.OrderUUID1
	expectedOrder := fixtures.ValidOrder()

	// Настраиваем мок
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(expectedOrder, nil)

	// Измеряем время выполнения множественных запросов
	const iterations = 1000
	start := time.Now()

	for i := 0; i < iterations; i++ {
		result, err := s.service.GetOrder(context.Background(), orderUUID)
		s.Require().NoError(err)
		s.Require().NotNil(result)
	}

	duration := time.Since(start)
	avgDuration := duration / iterations

	s.T().Logf("GetOrder: %d итераций за %v (среднее: %v на операцию)",
		iterations, duration, avgDuration)

	// Проверяем, что среднее время не превышает 1мс
	s.Require().Less(avgDuration, time.Millisecond,
		"GetOrder слишком медленный: %v", avgDuration)
}
