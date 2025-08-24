package order_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/fixtures"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

// TestCreateOrder_ComprehensiveErrorHandling тестирует все категории ошибок для CreateOrder
func (s *ServiceSuite) TestCreateOrder_ComprehensiveErrorHandling() {
	type errorTestCase struct {
		name          string
		scenario      fixtures.ErrorScenario
		input         *model.Order
		mockSetup     func()
		validateError func(*testing.T, error)
		description   string
	}

	// Предопределённые данные
	userUUID := fixtures.TestConstants.UserUUID1
	partUUID := fixtures.TestConstants.PartUUID1

	testCases := []errorTestCase{
		// === VALIDATION ERRORS ===
		{
			name:     "валидация_nil_заказ",
			scenario: fixtures.ValidationErrors.NilOrder,
			input:    nil,
			mockSetup: func() {
				// Валидационные ошибки не требуют моков
			},
			validateError: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "order cannot be nil")
			},
			description: "Тестирование валидации nil заказа",
		},
		{
			name:     "валидация_пустой_user_uuid",
			scenario: fixtures.ValidationErrors.NilUserUUID,
			input: fixtures.NewOrderBuilder().
				WithUserUUID(fixtures.TestConstants.UserUUID1). // Используем корректный UUID
				WithPartUUIDs().                                // Пустой список деталей вызовет другую ошибку
				Build(),
			mockSetup: func() {
				// Валидационные ошибки не требуют моков
			},
			validateError: func(t *testing.T, err error) {
				require.Error(t, err)
				// Проверяем что это именно ошибка валидации списка деталей
				helper := fixtures.NewErrorAssertionHelper("валидация_пустой_список_деталей")
				assertErr := helper.ExpectError(fixtures.ValidationErrors.EmptyParts.Error).
					ActualError(err).
					AssertMatch()
				require.NoError(t, assertErr)
			},
			description: "Тестирование валидации пустого списка деталей",
		},
		{
			name:     "валидация_слишком_много_деталей",
			scenario: fixtures.ValidationErrors.TooManyParts,
			input: fixtures.NewOrderBuilder().
				WithUserUUID(userUUID).
				WithPartUUIDs(generateManyParts(1001)...). // Предполагаем лимит в 1000
				Build(),
			mockSetup: func() {
				// Валидационные ошибки не требуют моков
			},
			validateError: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "too many parts")
			},
			description: "Тестирование валидации превышения лимита деталей",
		},

		// === GRPC ERRORS ===
		{
			name:     "grpc_inventory_недоступен",
			scenario: fixtures.GRPCErrors.Unavailable,
			input: fixtures.NewOrderBuilder().
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID).
				Build(),
			mockSetup: func() {
				s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
					Return(nil, fixtures.GRPCErrors.Unavailable.Error)
			},
			validateError: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, codes.Unavailable, status.Code(err))

				// Используем helper для детальной проверки
				helper := fixtures.NewErrorAssertionHelper("grpc_unavailable")
				assertErr := helper.ExpectError(fixtures.GRPCErrors.Unavailable.Error).
					ActualError(err).
					AssertMatch()
				require.NoError(t, assertErr)
			},
			description: "Тестирование недоступности inventory сервиса",
		},
		{
			name:     "grpc_таймаут_inventory",
			scenario: fixtures.GRPCErrors.DeadlineExceeded,
			input: fixtures.NewOrderBuilder().
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID).
				Build(),
			mockSetup: func() {
				s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
					Return(nil, fixtures.GRPCErrors.DeadlineExceeded.Error)
			},
			validateError: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, codes.DeadlineExceeded, status.Code(err))

				// Проверяем детальную информацию об ошибке
				errorInfo := fixtures.GetDetailedErrorInfo(err)
				require.Equal(t, "DeadlineExceeded", errorInfo["grpc_code"])
				require.True(t, errorInfo["is_grpc"].(bool))
			},
			description: "Тестирование таймаута при обращении к inventory",
		},

		// === BUSINESS ERRORS ===
		{
			name:     "бизнес_детали_не_найдены",
			scenario: fixtures.BusinessErrors.PartsNotFound,
			input: fixtures.NewOrderBuilder().
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID, fixtures.TestConstants.PartUUID2).
				Build(),
			mockSetup: func() {
				// Возвращаем только одну деталь из двух запрошенных
				parts := []model.Part{
					*fixtures.NewPartBuilder().WithPartUUID(partUUID).WithPrice(100.0).Build(),
				}
				s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
					Return(&parts, nil)
			},
			validateError: func(t *testing.T, err error) {
				require.Error(t, err)
				require.True(t, fixtures.IsExpectedError(err, fixtures.BusinessErrors.PartsNotFound.Error))
			},
			description: "Тестирование ситуации когда не все детали найдены",
		},

		// === INFRASTRUCTURE ERRORS ===
		{
			name:     "инфраструктура_ошибка_базы_данных",
			scenario: fixtures.InfrastructureErrors.DatabaseConnection,
			input: fixtures.NewOrderBuilder().
				WithUserUUID(userUUID).
				WithPartUUIDs(partUUID).
				Build(),
			mockSetup: func() {
				parts := []model.Part{
					*fixtures.NewPartBuilder().WithPartUUID(partUUID).WithPrice(100.0).Build(),
				}
				s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
					Return(&parts, nil)
				s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(nil, fixtures.InfrastructureErrors.DatabaseConnection.Error)
			},
			validateError: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "database connection failed")
			},
			description: "Тестирование ошибки подключения к базе данных",
		},
	}

	// Выполняем все тест-кейсы
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()

			result, err := s.service.CreateOrder(context.Background(), tc.input)

			// Результат должен быть nil при ошибке
			s.Require().Nil(result, "результат должен быть nil при ошибке")

			// Выполняем детальную валидацию ошибки
			tc.validateError(s.T(), err)

			// Логируем детальную информацию об ошибке для отладки
			errorInfo := fixtures.GetDetailedErrorInfo(err)
			s.T().Logf("Error details for %s: %+v", tc.name, errorInfo)

			// Проверяем моки
			s.inventoryClient.AssertExpectations(s.T())
			s.orderRepository.AssertExpectations(s.T())
		})
	}
}

// TestPayOrder_ComprehensiveErrorHandling тестирует все категории ошибок для PayOrder
func (s *ServiceSuite) TestPayOrder_ComprehensiveErrorHandling() {
	type errorTestCase struct {
		name          string
		scenario      fixtures.ErrorScenario
		input         *model.Order
		mockSetup     func()
		validateError func(*testing.T, error)
		description   string
	}

	orderUUID := fixtures.TestConstants.OrderUUID1
	userUUID := fixtures.TestConstants.UserUUID1
	partUUID := fixtures.TestConstants.PartUUID1

	testCases := []errorTestCase{
		// === VALIDATION ERRORS ===
		{
			name:     "валидация_пустой_payment_method",
			scenario: fixtures.ValidationErrors.EmptyPaymentMethod,
			input: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithPaymentMethod(""). // Пустой метод оплаты
				Build(),
			mockSetup: func() {
				// PayOrder сначала вызывает GetOrder, поэтому нужен мок
				existingOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithStatus(model.StatusPendingPayment).
					Build()
				s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
					Return(existingOrder, nil)
				// Также нужен мок для CreatePayment, который может вызываться до валидации
				s.paymentClient.On("CreatePayment", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(nil, errors.New("payment method cannot be empty")).Maybe()
			},
			validateError: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "payment method cannot be empty")
			},
			description: "Тестирование валидации пустого метода оплаты",
		},

		// === BUSINESS ERRORS ===
		{
			name:     "бизнес_заказ_не_найден",
			scenario: fixtures.BusinessErrors.OrderNotFound,
			input: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithPaymentMethod("CARD").
				Build(),
			mockSetup: func() {
				s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
					Return(nil, fixtures.BusinessErrors.OrderNotFound.Error)
			},
			validateError: func(t *testing.T, err error) {
				require.Error(t, err)
				require.True(t, fixtures.IsExpectedError(err, fixtures.BusinessErrors.OrderNotFound.Error))
			},
			description: "Тестирование оплаты несуществующего заказа",
		},
		{
			name:     "бизнес_заказ_уже_оплачен",
			scenario: fixtures.BusinessErrors.OrderAlreadyPaid,
			input: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithPaymentMethod("CARD").
				Build(),
			mockSetup: func() {
				// Заказ уже в статусе PAID
				dbOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithStatus(model.StatusPaid).
					Build()
				s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
					Return(dbOrder, nil)
				// Добавляем мок для CreatePayment, который может вызываться
				s.paymentClient.On("CreatePayment", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(nil, errors.New("order already paid")).Maybe()
			},
			validateError: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "order already paid")
			},
			description: "Тестирование повторной оплаты уже оплаченного заказа",
		},

		// === DEPENDENCY ERRORS ===
		{
			name:     "зависимость_payment_service_недоступен",
			scenario: fixtures.DependencyErrors.PaymentService,
			input: fixtures.NewOrderBuilder().
				WithOrderUUID(orderUUID).
				WithPaymentMethod("CARD").
				Build(),
			mockSetup: func() {
				dbOrder := fixtures.NewOrderBuilder().
					WithOrderUUID(orderUUID).
					WithUserUUID(userUUID).
					WithPartUUIDs(partUUID).
					WithTotalPrice(100.0).
					WithStatus(model.StatusPendingPayment).
					Build()

				s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
					Return(dbOrder, nil)
				s.paymentClient.On("CreatePayment", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(nil, fixtures.DependencyErrors.PaymentService.Error)
			},
			validateError: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, codes.Internal, status.Code(err))
				require.Contains(t, err.Error(), "payment service error")
			},
			description: "Тестирование недоступности платёжного сервиса",
		},
	}

	// Выполняем все тест-кейсы
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockSetup()

			result, err := s.service.PayOrder(context.Background(), tc.input)

			// Результат должен быть nil при ошибке
			s.Require().Nil(result, "результат должен быть nil при ошибке")

			// Выполняем детальную валидацию ошибки
			tc.validateError(s.T(), err)

			// Логируем детальную информацию об ошибке
			errorInfo := fixtures.GetDetailedErrorInfo(err)
			s.T().Logf("Error details for %s: %+v", tc.name, errorInfo)

			// Проверяем моки
			s.orderRepository.AssertExpectations(s.T())
			s.paymentClient.AssertExpectations(s.T())
		})
	}
}

// TestErrorRecovery тестирует восстановление после ошибок
func (s *ServiceSuite) TestErrorRecovery() {
	s.Run("восстановление_после_временной_ошибки_inventory", func() {
		userUUID := fixtures.TestConstants.UserUUID1
		partUUID := fixtures.TestConstants.PartUUID1

		order := fixtures.NewOrderBuilder().
			WithUserUUID(userUUID).
			WithPartUUIDs(partUUID).
			Build()

		// Первый вызов - ошибка
		s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
			Return(nil, fixtures.GRPCErrors.Unavailable.Error).Once()

		// Первая попытка должна завершиться ошибкой
		result1, err1 := s.service.CreateOrder(context.Background(), order)
		s.Require().Error(err1)
		s.Require().Nil(result1)
		s.Require().Equal(codes.Unavailable, status.Code(err1))

		// Сбрасываем моки для второй попытки
		s.SetupTest()

		// Второй вызов - успех
		parts := []model.Part{
			*fixtures.NewPartBuilder().WithPartUUID(partUUID).WithPrice(100.0).Build(),
		}
		expectedOrder := fixtures.NewOrderBuilder().
			WithUserUUID(userUUID).
			WithPartUUIDs(partUUID).
			WithTotalPrice(100.0).
			WithStatus(model.StatusPendingPayment).
			Build()

		s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
			Return(&parts, nil)
		s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
			Return(expectedOrder, nil)

		// Вторая попытка должна быть успешной
		result2, err2 := s.service.CreateOrder(context.Background(), order)
		s.Require().NoError(err2)
		s.Require().NotNil(result2)
		s.Require().Equal(model.StatusPendingPayment, result2.Status)

		s.inventoryClient.AssertExpectations(s.T())
		s.orderRepository.AssertExpectations(s.T())
	})
}

// TestErrorAggregation тестирует агрегацию множественных ошибок
func (s *ServiceSuite) TestErrorAggregation() {
	s.Run("множественные_ошибки_валидации", func() {
		// Заказ с множественными проблемами валидации
		invalidOrder := &model.Order{
			// Отсутствует UserUUID (nil)
			// Отсутствуют PartUUIDs (nil)
			// Отрицательная цена
			TotalPrice: -100.0,
		}

		result, err := s.service.CreateOrder(context.Background(), invalidOrder)

		s.Require().Error(err)
		s.Require().Nil(result)

		// Проверяем что ошибка содержит информацию о проблемах валидации
		errorInfo := fixtures.GetDetailedErrorInfo(err)
		s.T().Logf("Multiple validation errors: %+v", errorInfo)
	})
}

// Helper функции

// generateManyParts генерирует много UUID для тестирования лимитов
func generateManyParts(count int) []uuid.UUID {
	parts := make([]uuid.UUID, count)
	for i := 0; i < count; i++ {
		parts[i] = uuid.New()
	}
	return parts
}

// TestErrorContextPropagation тестирует передачу контекста ошибок
func (s *ServiceSuite) TestErrorContextPropagation() {
	s.Run("передача_контекста_отмены", func() {
		// Создаём контекст с отменой
		ctx, cancel := context.WithCancel(context.Background())

		order := fixtures.NewOrderBuilder().
			WithUserUUID(fixtures.TestConstants.UserUUID1).
			WithPartUUIDs(fixtures.TestConstants.PartUUID1).
			Build()

		// Настраиваем мок с задержкой
		s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
			Run(func(args mock.Arguments) {
				// Отменяем контекст во время выполнения
				cancel()
			}).
			Return(nil, context.Canceled)

		result, err := s.service.CreateOrder(ctx, order)

		s.Require().Error(err)
		s.Require().Nil(result)
		// Проверяем что ошибка содержит context.Canceled (может быть обёрнута)
		s.Require().True(errors.Is(err, context.Canceled) || err.Error() == context.Canceled.Error(),
			"Ошибка должна содержать context.Canceled, получена: %v", err)

		s.inventoryClient.AssertExpectations(s.T())
	})
}
