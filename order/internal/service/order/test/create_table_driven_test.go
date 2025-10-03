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

// TestCreateOrder_TableDriven использует table-driven подход для тестирования CreateOrder
// Это заменяет все отдельные тесты CreateOrder из create_test.go
func (s *ServiceSuite) TestCreateOrder_TableDriven() {
	// Определяем структуру тест-кейса
	type testCase struct {
		name            string
		input           *model.Order
		mockSetup       func()
		expectedResult  *model.Order
		expectedError   error
		validateResult  func(*testing.T, *model.Order) // дополнительная валидация результата
		validateMocks   func(*testing.T)               // валидация вызовов моков
		skipMockCleanup bool                           // пропустить очистку моков (для специальных случаев)
	}

	// Предопределенные UUID для консистентности тестов
	userUUID := fixtures.TestConstants.UserUUID1
	partUUID1 := fixtures.TestConstants.PartUUID1
	partUUID2 := fixtures.TestConstants.PartUUID2

	testCases := []testCase{
		{
			name: "успешное_создание_заказа_с_одной_деталью",
			input: fixtures.NewOrderBuilder().
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID1).
				Build(),
			mockSetup: func() {
				parts := []model.Part{
					*fixtures.NewPartBuilder().
						WithPartUUID(partUUID1).
						WithPrice(100.0).
						Build(),
				}
				expectedOrder := fixtures.NewOrderBuilder().
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID1).
					WithTotalPrice(100.0).
					WithStatus(model.StatusPendingPayment).
					Build()

				s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
					Return(&parts, nil)
				s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(expectedOrder, nil)
			},
			expectedResult: fixtures.NewOrderBuilder().
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID1).
				WithTotalPrice(100.0).
				WithStatus(model.StatusPendingPayment).
				Build(),
			expectedError: nil,
			validateResult: func(t *testing.T, result *model.Order) {
				require.Equal(t, model.StatusPendingPayment, result.Status)
				require.Equal(t, float32(100.0), result.TotalPrice)
				require.Len(t, result.PartUUIDs, 1)
			},
		},
		{
			name: "успешное_создание_заказа_с_несколькими_деталями",
			input: fixtures.NewOrderBuilder().
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID1, partUUID2).
				Build(),
			mockSetup: func() {
				parts := []model.Part{
					*fixtures.NewPartBuilder().WithPartUUID(partUUID1).WithPrice(100.0).Build(),
					*fixtures.NewPartBuilder().WithPartUUID(partUUID2).WithPrice(200.0).Build(),
				}
				expectedOrder := fixtures.NewOrderBuilder().
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID1, partUUID2).
					WithTotalPrice(300.0).
					WithStatus(model.StatusPendingPayment).
					Build()

				s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
					Return(&parts, nil)
				s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(expectedOrder, nil)
			},
			expectedResult: fixtures.NewOrderBuilder().
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID1, partUUID2).
				WithTotalPrice(300.0).
				WithStatus(model.StatusPendingPayment).
				Build(),
			expectedError: nil,
			validateResult: func(t *testing.T, result *model.Order) {
				require.Equal(t, model.StatusPendingPayment, result.Status)
				require.Equal(t, float32(300.0), result.TotalPrice)
				require.Len(t, result.PartUUIDs, 2)
			},
		},
		{
			name:  "ошибка_валидации_пустой_список_деталей",
			input: fixtures.OrderWithEmptyParts(),
			mockSetup: func() {
				// Моки не нужны - валидация происходит до обращения к зависимостям
			},
			expectedResult: nil,
			expectedError:  model.ErrPartsSpecified,
		},
		{
			name:  "ошибка_валидации_nil_заказа",
			input: nil,
			mockSetup: func() {
				// Моки не нужны - валидация происходит до обращения к зависимостям
			},
			expectedResult: nil,
			expectedError:  errors.New("order cannot be nil"),
		},
		{
			name: "ошибка_валидации_nil_user_uuid",
			input: fixtures.NewOrderBuilder().
				WithUserUUID(uuid.Nil).
				WithPartUUIDs(partUUID1).
				Build(),
			mockSetup: func() {
				// Моки не нужны - валидация происходит до обращения к зависимостям
			},
			expectedResult: nil,
			expectedError:  fixtures.ValidationErrors.NilUserUUID.Error,
			validateResult: func(t *testing.T, result *model.Order) {
				require.Nil(t, result)
			},
		},
		{
			name: "частичные_данные_от_inventory_service",
			input: fixtures.NewOrderBuilder().
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID1, partUUID2).
				Build(),
			mockSetup: func() {
				// Возвращаем только одну деталь из двух запрошенных
				parts := []model.Part{
					*fixtures.NewPartBuilder().WithPartUUID(partUUID1).WithPrice(100.0).Build(),
				}
				s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
					Return(&parts, nil)
			},
			expectedResult: nil,
			expectedError:  model.ErrPartsListNotFound,
			validateMocks: func(t *testing.T) {
				// Проверяем, что repository.CreateOrder НЕ вызывался
				s.orderRepository.AssertNotCalled(t, "CreateOrder")
			},
		},
		{
			name: "ошибка_inventory_service_недоступен",
			input: fixtures.NewOrderBuilder().
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID1).
				Build(),
			mockSetup: func() {
				s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
					Return(nil, fixtures.GRPCErrors.Unavailable.Error)
			},
			expectedResult: nil,
			expectedError:  fixtures.GRPCErrors.Unavailable.Error,
			validateResult: func(t *testing.T, result *model.Order) {
				require.Nil(t, result)
			},
		},
		{
			name: "ошибка_inventory_service_таймаут",
			input: fixtures.NewOrderBuilder().
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID1).
				Build(),
			mockSetup: func() {
				s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
					Return(nil, fixtures.GRPCErrors.DeadlineExceeded.Error)
			},
			expectedResult: nil,
			expectedError:  fixtures.GRPCErrors.DeadlineExceeded.Error,
			validateResult: func(t *testing.T, result *model.Order) {
				require.Nil(t, result)
			},
		},
		{
			name: "ошибка_создания_в_repository",
			input: fixtures.NewOrderBuilder().
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID1, partUUID2).
				Build(),
			mockSetup: func() {
				parts := []model.Part{
					*fixtures.NewPartBuilder().WithPartUUID(partUUID1).WithPrice(100.0).Build(),
					*fixtures.NewPartBuilder().WithPartUUID(partUUID2).WithPrice(200.0).Build(),
				}
				s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
					Return(&parts, nil)
				s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(nil, model.ErrOrderNotFound)
			},
			expectedResult: nil,
			expectedError:  model.ErrOrderNotFound,
		},
		{
			name: "граничный_случай_максимальная_цена",
			input: fixtures.NewOrderBuilder().
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID1).
				Build(),
			mockSetup: func() {
				parts := []model.Part{
					*fixtures.NewPartBuilder().
						WithPartUUID(partUUID1).
						WithPrice(999999999.99).
						Build(),
				}
				expectedOrder := fixtures.NewOrderBuilder().
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID1).
					WithTotalPrice(999999999.99).
					WithStatus(model.StatusPendingPayment).
					Build()

				s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
					Return(&parts, nil)
				s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(expectedOrder, nil)
			},
			expectedResult: fixtures.NewOrderBuilder().
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID1).
				WithTotalPrice(999999999.99).
				WithStatus(model.StatusPendingPayment).
				Build(),
			expectedError: nil,
			validateResult: func(t *testing.T, result *model.Order) {
				require.True(t, result.TotalPrice >= 999999999.0, "Цена должна быть близка к максимальной")
				require.Equal(t, model.StatusPendingPayment, result.Status)
			},
		},
		{
			name: "граничный_случай_минимальная_цена",
			input: fixtures.NewOrderBuilder().
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID1).
				Build(),
			mockSetup: func() {
				parts := []model.Part{
					*fixtures.NewPartBuilder().
						WithPartUUID(partUUID1).
						WithPrice(0.01).
						Build(),
				}
				expectedOrder := fixtures.NewOrderBuilder().
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID1).
					WithTotalPrice(0.01).
					WithStatus(model.StatusPendingPayment).
					Build()

				s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
					Return(&parts, nil)
				s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(expectedOrder, nil)
			},
			expectedResult: fixtures.NewOrderBuilder().
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID1).
				WithTotalPrice(0.01).
				WithStatus(model.StatusPendingPayment).
				Build(),
			expectedError: nil,
		},
	}

	// Выполняем все тест-кейсы
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Очищаем моки перед каждым тестом
			if !tc.skipMockCleanup {
				s.SetupTest()
			}

			// Настраиваем моки для текущего тест-кейса
			tc.mockSetup()

			// Вызываем тестируемый метод
			result, err := s.service.CreateOrder(context.Background(), tc.input)

			// Проверяем ошибку с использованием улучшенного error handling
			if tc.expectedError != nil {
				s.Require().Error(err, "ожидалась ошибка для сценария: %s", tc.name)

				// Используем централизованную проверку ошибок
				isMatch := fixtures.IsExpectedError(err, tc.expectedError)
				s.Require().True(isMatch, "ошибка не соответствует ожидаемой. Ожидалась: %v, получена: %v", tc.expectedError, err)

				s.Require().Nil(result, "результат должен быть nil при ошибке")

				// Логируем детальную информацию об ошибке для отладки
				errorInfo := fixtures.GetDetailedErrorInfo(err)
				s.T().Logf("Error details for %s: %+v", tc.name, errorInfo)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(result)

				// Проверяем основные поля результата
				s.Require().Equal(tc.expectedResult.UserUUID, result.UserUUID)
				s.Require().Equal(float32(tc.expectedResult.TotalPrice), result.TotalPrice) // Приводим к float32
				s.Require().Equal(tc.expectedResult.Status, result.Status)
				s.Require().Equal(len(tc.expectedResult.PartUUIDs), len(result.PartUUIDs))

				// Выполняем дополнительную валидацию, если она определена
				if tc.validateResult != nil {
					tc.validateResult(s.T(), result)
				}
			}

			// Выполняем дополнительную валидацию моков, если определена
			if tc.validateMocks != nil {
				tc.validateMocks(s.T())
			} else {
				// Стандартная проверка моков
				s.inventoryClient.AssertExpectations(s.T())
				s.orderRepository.AssertExpectations(s.T())
			}
		})
	}
}

// TestCreateOrder_ConcurrentSafety тестирует thread-safety CreateOrder
func (s *ServiceSuite) TestCreateOrder_ConcurrentSafety() {
	const numGoroutines = 10

	testCases := []struct {
		name  string
		setup func() (*model.Order, []model.Part, *model.Order)
	}{
		{
			name: "параллельное_создание_разных_заказов",
			setup: func() (*model.Order, []model.Part, *model.Order) {
				userUUID := uuid.New()
				partUUID := uuid.New()

				inputOrder := fixtures.NewOrderBuilder().
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					Build()

				parts := []model.Part{
					*fixtures.NewPartBuilder().
						WithPartUUID(partUUID).
						WithPrice(100.0).
						Build(),
				}

				expectedOrder := fixtures.NewOrderBuilder().
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(100.0).
					WithStatus(model.StatusPendingPayment).
					Build()

				return inputOrder, parts, expectedOrder
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Убираем t.Parallel() для избежания конфликтов

			results := make(chan error, numGoroutines)

			// Настраиваем моки для всех горутин
			for i := 0; i < numGoroutines; i++ {
				_, parts, expectedOrder := tc.setup()

				s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
					Return(&parts, nil).Once()
				s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(expectedOrder, nil).Once()
			}

			// Запускаем горутины
			for i := 0; i < numGoroutines; i++ {
				go func() {
					order, _, _ := tc.setup()
					_, err := s.service.CreateOrder(context.Background(), order)
					results <- err
				}()
			}

			// Собираем результаты
			for i := 0; i < numGoroutines; i++ {
				select {
				case err := <-results:
					s.Require().NoError(err)
				case <-context.Background().Done():
					s.Fail("тест завис")
				}
			}

			s.inventoryClient.AssertExpectations(s.T())
			s.orderRepository.AssertExpectations(s.T())
		})
	}
}

// TestCreateOrder_InputValidation тестирует валидацию входных параметров
func (s *ServiceSuite) TestCreateOrder_InputValidation() {
	// Убираем t.Parallel() для избежания конфликтов

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
			name: "пустой_user_uuid",
			order: fixtures.NewOrderBuilder().
				WithUserUUID(uuid.Nil).
				Build(),
			expectedError: "user UUID cannot be nil",
		},
		{
			name:          "пустой_список_деталей",
			order:         fixtures.OrderWithEmptyParts(),
			expectedError: "parts not specified",
		},
	}

	for _, vc := range validationCases {
		vc := vc // захватываем переменную для goroutine
		s.Run(vc.name, func() {
			// Убираем t.Parallel() для избежания конфликтов
			s.SetupTest()

			// Для валидационных тестов моки не нужны
			_, err := s.service.CreateOrder(context.Background(), vc.order)

			s.Require().Error(err)
			s.Require().Contains(err.Error(), vc.expectedError)
		})
	}
}
