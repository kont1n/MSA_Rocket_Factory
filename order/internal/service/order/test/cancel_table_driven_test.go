package order_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/fixtures"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

// TestCancelOrder_TableDriven использует table-driven подход для тестирования CancelOrder
// Это заменяет все отдельные тесты CancelOrder из cancel_test.go
func (s *ServiceSuite) TestCancelOrder_TableDriven() {
	type testCase struct {
		name           string
		input          *model.Order
		mockSetup      func(*model.Order) // функция получает input заказ для настройки моков
		expectedResult *model.Order
		expectedError  error
		validateResult func(*testing.T, *model.Order)
		description    string
	}

	// Предопределенные данные
	orderUUID := fixtures.TestConstants.OrderUUID1
	userUUID := fixtures.TestConstants.UserUUID1
	partUUID := fixtures.TestConstants.PartUUID1

	testCases := []testCase{
		{
			name: "успешная_отмена_заказа_в_ожидании_оплаты",
			input: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID).
				WithTotalPrice(100.0).
				WithStatus(model.StatusPendingPayment).
				Build(),
			mockSetup: func(inputOrder *model.Order) {
				// Заказ существует в БД в статусе pending payment
				dbOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(inputOrder.OrderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(100.0).
					WithStatus(model.StatusPendingPayment).
					Build()

				// Заказ после отмены
				cancelledOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(inputOrder.OrderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(100.0).
					WithStatus(model.StatusCancelled).
					Build()

				s.orderRepository.On("GetOrder", mock.Anything, inputOrder.OrderUUID).
					Return(dbOrder, nil)
				s.orderRepository.On("UpdateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(cancelledOrder, nil)
			},
			expectedResult: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID).
				WithTotalPrice(100.0).
				WithStatus(model.StatusCancelled).
				Build(),
			expectedError: nil,
			validateResult: func(t *testing.T, result *model.Order) {
				require.Equal(t, model.StatusCancelled, result.Status)
				require.Equal(t, orderUUID, result.OrderUUID)
				require.Equal(t, userUUID, result.UserUUID)
			},
			description: "Стандартный сценарий отмены заказа, ожидающего оплаты",
		},
		{
			name: "ошибка_отмена_заказа_оплаченного_не_разрешена",
			input: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID).
				WithStatus(model.StatusPaid).
				Build(),
			mockSetup: func(inputOrder *model.Order) {
				// Заказ оплачен
				dbOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(inputOrder.OrderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(150.0).
					WithStatus(model.StatusPaid).
					Build()

				cancelledOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(inputOrder.OrderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(150.0).
					WithStatus(model.StatusCancelled).
					Build()

				s.orderRepository.On("GetOrder", mock.Anything, inputOrder.OrderUUID).
					Return(dbOrder, nil)
				s.orderRepository.On("UpdateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(cancelledOrder, nil)
			},
			expectedResult: nil,
			expectedError:  model.ErrPaid,
			description:    "Попытка отмены оплаченного заказа (не разрешено бизнес-логикой)",
		},
		{
			name: "ошибка_заказ_не_найден",
			input: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				Build(),
			mockSetup: func(inputOrder *model.Order) {
				s.orderRepository.On("GetOrder", mock.Anything, inputOrder.OrderUUID).
					Return(nil, model.ErrOrderNotFound)
			},
			expectedResult: nil,
			expectedError:  model.ErrOrderNotFound,
			description:    "Попытка отмены несуществующего заказа",
		},
		{
			name: "ошибка_заказ_уже_оплачен",
			input: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithStatus(model.StatusPendingPayment).
				Build(),
			mockSetup: func(inputOrder *model.Order) {
				// Заказ уже оплачен в БД
				dbOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(inputOrder.OrderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(200.0).
					WithStatus(model.StatusPaid).
					Build()

				s.orderRepository.On("GetOrder", mock.Anything, inputOrder.OrderUUID).
					Return(dbOrder, nil)
			},
			expectedResult: nil,
			expectedError:  model.ErrPaid,
			description:    "Попытка отмены уже оплаченного заказа",
		},
		{
			name: "ошибка_заказ_уже_отменён",
			input: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithStatus(model.StatusPendingPayment).
				Build(),
			mockSetup: func(inputOrder *model.Order) {
				// Заказ уже отменён в БД
				dbOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(inputOrder.OrderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(100.0).
					WithStatus(model.StatusCancelled).
					Build()

				s.orderRepository.On("GetOrder", mock.Anything, inputOrder.OrderUUID).
					Return(dbOrder, nil)
			},
			expectedResult: nil,
			expectedError:  model.ErrCancelled,
			description:    "Попытка отмены уже отменённого заказа",
		},
		// Удаляем тест с несуществующим статусом StatusAssembled
		{
			name: "ошибка_обновления_в_репозитории",
			input: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithUserUUID(userUUID).
				WithStatus(model.StatusPendingPayment).
				Build(),
			mockSetup: func(inputOrder *model.Order) {
				// Заказ существует и может быть отменён
				dbOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(inputOrder.OrderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(100.0).
					WithStatus(model.StatusPendingPayment).
					Build()

				s.orderRepository.On("GetOrder", mock.Anything, inputOrder.OrderUUID).
					Return(dbOrder, nil)
				// Но обновление не проходит
				s.orderRepository.On("UpdateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(nil, model.ErrOrderNotFound)
			},
			expectedResult: nil,
			expectedError:  model.ErrOrderNotFound,
			description:    "Ошибка обновления заказа в репозитории",
		},
		{
			name: "граничный_случай_заказ_с_большой_суммой",
			input: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithUserUUID(userUUID).
				WithTotalPrice(9999999.99).
				WithStatus(model.StatusPendingPayment).
				Build(),
			mockSetup: func(inputOrder *model.Order) {
				dbOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(inputOrder.OrderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(9999999.99).
					WithStatus(model.StatusPendingPayment).
					Build()

				cancelledOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(inputOrder.OrderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(9999999.99).
					WithStatus(model.StatusCancelled).
					Build()

				s.orderRepository.On("GetOrder", mock.Anything, inputOrder.OrderUUID).
					Return(dbOrder, nil)
				s.orderRepository.On("UpdateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(cancelledOrder, nil)
			},
			expectedResult: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID).
				WithTotalPrice(9999999.99).
				WithStatus(model.StatusCancelled).
				Build(),
			expectedError: nil,
			validateResult: func(t *testing.T, result *model.Order) {
				require.True(t, result.TotalPrice > float32(9999999.0))
				require.Equal(t, model.StatusCancelled, result.Status)
			},
			description: "Отмена заказа с максимальной суммой",
		},
	}

	// Выполняем все тест-кейсы
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Очищаем моки перед каждым тестом
			s.SetupTest()

			// Настраиваем моки для текущего тест-кейса
			tc.mockSetup(tc.input)

			// Вызываем тестируемый метод
			result, err := s.service.CancelOrder(context.Background(), tc.input)

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
				s.Require().Equal(tc.expectedResult.UserUUID, result.UserUUID)
				s.Require().Equal(tc.expectedResult.TotalPrice, result.TotalPrice)

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

// TestCancelOrder_StatusTransitions тестирует все возможные переходы статусов при отмене
func (s *ServiceSuite) TestCancelOrder_StatusTransitions() {
	statusTransitionCases := []struct {
		currentStatus model.OrderStatus
		canCancel     bool
		expectedError error
		description   string
	}{
		{
			currentStatus: model.StatusPendingPayment,
			canCancel:     true,
			expectedError: nil,
			description:   "Можно отменить заказ в ожидании оплаты",
		},
		// Удаляем тест с несуществующим статусом StatusAssembling
		{
			currentStatus: model.StatusPaid,
			canCancel:     false,
			expectedError: model.ErrPaid,
			description:   "Нельзя отменить оплаченный заказ",
		},
		{
			currentStatus: model.StatusCancelled,
			canCancel:     false,
			expectedError: model.ErrCancelled,
			description:   "Нельзя отменить уже отменённый заказ",
		},
		// Удаляем тест с несуществующим статусом StatusAssembled
	}

	orderUUID := fixtures.TestConstants.OrderUUID1
	userUUID := fixtures.TestConstants.UserUUID1
	partUUID := fixtures.TestConstants.PartUUID1

	for _, tc := range statusTransitionCases {
		s.Run("переход_из_"+string(tc.currentStatus), func() {
			s.SetupTest()

			// Входящий заказ
			inputOrder := fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID).
				WithStatus(model.StatusPendingPayment). // Статус в запросе не важен
				Build()

			// Заказ в БД с текущим статусом
			dbOrder := fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID).
				WithTotalPrice(100.0).
				WithStatus(tc.currentStatus).
				Build()

			s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
				Return(dbOrder, nil)

			if tc.canCancel {
				// Если отмена возможна, настраиваем мок для UpdateOrder
				cancelledOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(100.0).
					WithStatus(model.StatusCancelled).
					Build()

				s.orderRepository.On("UpdateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(cancelledOrder, nil)
			}

			// Вызываем метод
			result, err := s.service.CancelOrder(context.Background(), inputOrder)

			// Проверяем результат
			if tc.canCancel {
				s.Require().NoError(err, tc.description)
				s.Require().NotNil(result)
				s.Require().Equal(model.StatusCancelled, result.Status)
			} else {
				s.Require().Error(err, tc.description)
				s.Require().ErrorIs(err, tc.expectedError)
				s.Require().Nil(result)
			}

			s.orderRepository.AssertExpectations(s.T())
		})
	}
}

// TestCancelOrder_ConcurrentCancellation тестирует конкурентную отмену заказов
func (s *ServiceSuite) TestCancelOrder_ConcurrentCancellation() {
	// Убираем t.Parallel() для избежания конфликтов с параллельным выполнением

	const numGoroutines = 5
	orderUUID := fixtures.TestConstants.OrderUUID1
	userUUID := fixtures.TestConstants.UserUUID1
	partUUID := fixtures.TestConstants.PartUUID1

	// Подготавливаем данные
	inputOrder := fixtures.NewOrderBuilder().
		WithOrderUUID(orderUUID).
		WithUserUUID(userUUID).
		WithPartUUIDs(partUUID).
		WithStatus(model.StatusPendingPayment).
		Build()

	dbOrder := fixtures.NewOrderBuilder().
		WithOrderUUID(orderUUID).
		WithUserUUID(userUUID).
		WithPartUUIDs(partUUID).
		WithTotalPrice(100.0).
		WithStatus(model.StatusPendingPayment).
		Build()

	cancelledOrder := fixtures.NewOrderBuilder().
		WithOrderUUID(orderUUID).
		WithUserUUID(userUUID).
		WithPartUUIDs(partUUID).
		WithTotalPrice(100.0).
		WithStatus(model.StatusCancelled).
		Build()

	// Настраиваем моки для параллельных запросов
	// Только первый запрос должен успешно отменить заказ
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(dbOrder, nil).Once()
	s.orderRepository.On("UpdateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
		Return(cancelledOrder, nil).Once()

	// Остальные запросы должны получить уже отменённый заказ
	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
		Return(cancelledOrder, nil).Times(numGoroutines - 1)

	results := make(chan *model.Order, numGoroutines)
	errChan := make(chan error, numGoroutines)

	// Запускаем горутины
	for i := 0; i < numGoroutines; i++ {
		go func() {
			result, err := s.service.CancelOrder(context.Background(), inputOrder)
			results <- result
			errChan <- err
		}()
	}

	// Собираем результаты
	successCount := 0
	cancelledErrorCount := 0

	for i := 0; i < numGoroutines; i++ {
		result := <-results
		err := <-errChan

		switch {
		case err == nil:
			successCount++
			s.Require().NotNil(result)
			s.Require().Equal(model.StatusCancelled, result.Status)
		case errors.Is(err, model.ErrCancelled):
			cancelledErrorCount++
			s.Require().Nil(result)
		default:
			s.Failf("Неожиданная ошибка", "получена неожиданная ошибка: %v", err)
		}
	}

	// Должно быть ровно одно успешное выполнение
	s.Require().Equal(1, successCount, "должно быть ровно одно успешное выполнение")
	s.Require().Equal(numGoroutines-1, cancelledErrorCount, "остальные запросы должны получить ошибку ErrCancelled")

	s.orderRepository.AssertExpectations(s.T())
}

// TestCancelOrder_InputValidation тестирует валидацию входных параметров
func (s *ServiceSuite) TestCancelOrder_InputValidation() {
	validationCases := []struct {
		name          string
		order         *model.Order
		expectedError string
	}{
		{
			name:          "nil_заказ",
			order:         nil,
			expectedError: "order cannot be nil",
		},
		{
			name: "пустой_order_uuid",
			order: fixtures.NewOrderBuilder().
				WithOrderUUID(uuid.Nil).
				Build(),
			expectedError: "order UUID cannot be nil",
		},
	}

	for _, vc := range validationCases {
		vc := vc // захватываем переменную для goroutine
		s.Run(vc.name, func() {
			// Убираем дублирующий вызов t.Parallel()

			// Для валидационных тестов моки не нужны
			result, err := s.service.CancelOrder(context.Background(), vc.order)

			s.Require().Error(err)
			s.Require().Contains(err.Error(), vc.expectedError)
			s.Require().Nil(result)
		})
	}
}
