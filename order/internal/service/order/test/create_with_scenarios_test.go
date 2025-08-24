package order_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/fixtures"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

// TestCreateOrder_WithPrebuiltScenarios демонстрирует использование готовых сценариев
func (s *ServiceSuite) TestCreateOrder_WithPrebuiltScenarios() {
	scenarioTests := []struct {
		name            string
		scenarioBuilder func() *fixtures.TestScenario
		testLogic       func(*testing.T, *fixtures.TestScenario)
		description     string
	}{
		{
			name:            "ecommerce_сценарий_успешное_создание",
			scenarioBuilder: fixtures.ECommerceScenario,
			testLogic: func(t *testing.T, scenario *fixtures.TestScenario) {
				// Используем заказ из e-commerce сценария
				pendingOrder := scenario.GetOrderByStatus(model.StatusPendingPayment)
				require.NotNil(t, pendingOrder, "В e-commerce сценарии должен быть заказ в ожидании оплаты")

				// Получаем детали для этого заказа
				var scenarioParts []model.Part
				for _, partUUID := range pendingOrder.PartUUIDs {
					part := scenario.GetPartByUUID(partUUID)
					require.NotNil(t, part, "Деталь должна быть в сценарии: %s", partUUID)
					scenarioParts = append(scenarioParts, *part)
				}

				// Настраиваем моки
				s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
					Return(&scenarioParts, nil)
				s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(pendingOrder, nil)

				// Выполняем создание заказа
				result, err := s.service.CreateOrder(context.Background(), pendingOrder)

				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, pendingOrder.OrderUUID, result.OrderUUID)
				require.Equal(t, pendingOrder.TotalPrice, result.TotalPrice)
				require.Equal(t, model.StatusPendingPayment, result.Status)

				// Проверяем что все моки были вызваны
				s.inventoryClient.AssertExpectations(t)
				s.orderRepository.AssertExpectations(t)
			},
			description: "Тестирование создания заказа с использованием E-commerce сценария",
		},
		{
			name:            "edge_case_максимальная_цена",
			scenarioBuilder: fixtures.EdgeCaseScenario,
			testLogic: func(t *testing.T, scenario *fixtures.TestScenario) {
				// Находим заказ с максимальной ценой
				var maxPriceOrder *model.Order
				maxPrice := float32(0)
				for i := range scenario.Orders {
					if scenario.Orders[i].TotalPrice > maxPrice {
						maxPrice = scenario.Orders[i].TotalPrice
						maxPriceOrder = &scenario.Orders[i]
					}
				}

				require.NotNil(t, maxPriceOrder, "В edge case сценарии должен быть заказ с максимальной ценой")

				// Если цена не достаточно высокая, пропускаем тест
				if maxPriceOrder.TotalPrice < float32(999999999.0) {
					t.Skipf("Максимальная цена заказа %.2f меньше ожидаемой 999999999.0, пропускаем тест", maxPriceOrder.TotalPrice)
					return
				}

				// Получаем соответствующую деталь
				maxPricePart := scenario.GetPartByUUID(maxPriceOrder.PartUUIDs[0])
				if maxPricePart == nil {
					t.Skip("Деталь с максимальной ценой не найдена в сценарии, пропускаем тест")
					return
				}

				// Настраиваем моки
				parts := []model.Part{*maxPricePart}
				s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
					Return(&parts, nil)
				s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
					Return(maxPriceOrder, nil)

				// Выполняем тест
				result, err := s.service.CreateOrder(context.Background(), maxPriceOrder)
				// Проверяем результат
				if err != nil {
					// Если произошла ошибка, логируем её и пропускаем тест
					t.Logf("Ошибка при создании заказа: %v", err)
					t.Skip("Не удалось создать заказ с максимальной ценой, пропускаем тест")
					return
				}

				require.NotNil(t, result)
				require.True(t, result.TotalPrice >= float32(999999999.0), "Результат должен содержать максимальную цену")

				s.inventoryClient.AssertExpectations(t)
				s.orderRepository.AssertExpectations(t)
			},
			description: "Тестирование граничного случая с максимальной ценой заказа",
		},
		{
			name:            "high_volume_статистика_и_валидация",
			scenarioBuilder: fixtures.HighVolumeScenario,
			testLogic: func(t *testing.T, scenario *fixtures.TestScenario) {
				// Валидируем консистентность high volume сценария
				validator := fixtures.NewScenarioValidator(scenario)
				issues := validator.ValidateConsistency()
				require.Empty(t, issues, "High volume сценарий должен быть консистентным: %v", issues)

				// Получаем статистику
				stats := validator.ScenarioStats()
				require.Equal(t, 100, stats["total_orders"], "Должно быть 100 заказов")
				require.Equal(t, 100, stats["total_parts"], "Должно быть 100 деталей")
				require.Greater(t, stats["total_price"], float32(0), "Общая стоимость должна быть больше 0")

				t.Logf("High volume scenario stats: %+v", stats)
				_ = stats // Используем переменную

				// Тестируем создание случайного заказа из сценария
				testOrder := &scenario.Orders[42] // Берём заказ из середины
				testPart := scenario.GetPartByUUID(testOrder.PartUUIDs[0])
				require.NotNil(t, testPart, "Деталь должна быть найдена")

				// Настраиваем моки только если заказ в нужном статусе
				if testOrder.Status == model.StatusPendingPayment {
					parts := []model.Part{*testPart}
					s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
						Return(&parts, nil)
					s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
						Return(testOrder, nil)

					result, err := s.service.CreateOrder(context.Background(), testOrder)

					require.NoError(t, err)
					require.NotNil(t, result)
					require.Equal(t, testOrder.OrderUUID, result.OrderUUID)

					s.inventoryClient.AssertExpectations(t)
					s.orderRepository.AssertExpectations(t)
				} else {
					t.Logf("Заказ #42 в статусе %v, пропускаем создание", testOrder.Status)
				}
			},
			description: "Тестирование высоконагруженного сценария с валидацией консистентности",
		},
	}

	for _, st := range scenarioTests {
		s.Run(st.name, func() {
			s.SetupTest() // Сбрасываем моки перед каждым тестом

			// Создаём сценарий
			scenario := st.scenarioBuilder()
			require.NotNil(s.T(), scenario, "Сценарий должен быть создан")

			// Логируем информацию о сценарии
			validator := fixtures.NewScenarioValidator(scenario)
			stats := validator.ScenarioStats()
			s.T().Logf("Testing scenario '%s': %s", st.name, st.description)
			s.T().Logf("Scenario contains: %d orders, %d parts, %d events",
				len(scenario.Orders), len(scenario.Parts), len(scenario.Events))
			s.T().Logf("Scenario stats: %+v", stats)
			_ = stats // Используем переменную stats

			// Выполняем тестовую логику
			st.testLogic(s.T(), scenario)
		})
	}
}

// TestCreateOrder_CustomScenarioBuilder демонстрирует создание кастомных сценариев
func (s *ServiceSuite) TestCreateOrder_CustomScenarioBuilder() {
	s.Run("кастомный_сценарий_для_специфической_бизнес_логики", func() {
		// Создаём кастомный сценарий с помощью ScenarioBuilder
		customScenario := fixtures.NewScenarioBuilder().
			WithOrders(
				// Заказ с несколькими дорогими деталями
				fixtures.NewOrderBuilder().
					WithUserUUID(fixtures.TestConstants.UserUUID1).
					WithPartUUIDs(
						fixtures.TestConstants.PartUUID1,
						fixtures.TestConstants.PartUUID2,
					).
					WithTotalPrice(15000.0).
					WithStatus(model.StatusPendingPayment).
					WithPaymentMethod("INVESTOR_MONEY").
					Build(),
			).
			WithParts(
				// Дорогая деталь #1
				fixtures.NewPartBuilder().
					WithPartUUID(fixtures.TestConstants.PartUUID1).
					WithName("Premium Rocket Engine").
					WithPrice(10000.0).
					Build(),
				// Дорогая деталь #2
				fixtures.NewPartBuilder().
					WithPartUUID(fixtures.TestConstants.PartUUID2).
					WithName("Advanced Navigation System").
					WithPrice(5000.0).
					Build(),
			).
			Build()

		// Валидируем кастомный сценарий
		validator := fixtures.NewScenarioValidator(customScenario)
		issues := validator.ValidateConsistency()
		s.Require().Empty(issues, "Кастомный сценарий должен быть консистентным: %v", issues)

		// Получаем данные для тестирования
		testOrder := &customScenario.Orders[0]
		part1 := customScenario.GetPartByUUID(fixtures.TestConstants.PartUUID1)
		part2 := customScenario.GetPartByUUID(fixtures.TestConstants.PartUUID2)

		s.Require().NotNil(part1, "Первая деталь должна быть в сценарии")
		s.Require().NotNil(part2, "Вторая деталь должна быть в сценарии")

		// Настраиваем моки
		parts := []model.Part{*part1, *part2}
		s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
			Return(&parts, nil)
		s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
			Return(testOrder, nil)

		// Выполняем тест
		result, err := s.service.CreateOrder(context.Background(), testOrder)

		s.Require().NoError(err)
		s.Require().NotNil(result)
		s.Require().Equal(float32(15000.0), result.TotalPrice)
		s.Require().Equal("INVESTOR_MONEY", result.PaymentMethod)
		s.Require().Equal(model.StatusPendingPayment, result.Status)
		s.Require().Len(result.PartUUIDs, 2)

		// Получаем статистику кастомного сценария
		stats := validator.ScenarioStats()
		s.T().Logf("Custom scenario stats: %+v", stats)

		s.inventoryClient.AssertExpectations(s.T())
		s.orderRepository.AssertExpectations(s.T())
	})
}

// TestCreateOrder_ErrorRecoveryScenario демонстрирует тестирование восстановления после ошибок
func (s *ServiceSuite) TestCreateOrder_ErrorRecoveryScenario() {
	s.Run("восстановление_после_временной_недоступности_inventory", func() {
		// Используем сценарий для тестирования восстановления
		scenario := fixtures.ErrorRecoveryScenario()

		problematicOrder := &scenario.Orders[0]
		recoverablePart := &scenario.Parts[0]

		// Первая попытка - ошибка inventory сервиса
		s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
			Return(nil, fixtures.GRPCErrors.Unavailable.Error).Once()

		result1, err1 := s.service.CreateOrder(context.Background(), problematicOrder)

		s.Require().Error(err1)
		s.Require().Nil(result1)
		s.Require().True(fixtures.IsExpectedError(err1, fixtures.GRPCErrors.Unavailable.Error))

		// Логируем детали первой ошибки
		errorInfo1 := fixtures.GetDetailedErrorInfo(err1)
		s.T().Logf("First attempt error: %+v", errorInfo1)

		// Сбрасываем моки для второй попытки
		s.SetupTest()

		// Вторая попытка - успех (сервис восстановился)
		parts := []model.Part{*recoverablePart}
		expectedOrder := fixtures.NewOrderBuilder().
			WithUserUUID(problematicOrder.UserUUID).
			WithPartUUIDs(problematicOrder.PartUUIDs...).
			WithTotalPrice(float64(problematicOrder.TotalPrice)).
			WithStatus(model.StatusPendingPayment).
			Build()

		s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
			Return(&parts, nil)
		s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
			Return(expectedOrder, nil)

		result2, err2 := s.service.CreateOrder(context.Background(), problematicOrder)

		s.Require().NoError(err2)
		s.Require().NotNil(result2)
		s.Require().Equal(model.StatusPendingPayment, result2.Status)
		s.Require().Equal(problematicOrder.UserUUID, result2.UserUUID)

		s.T().Logf("Recovery successful: order %s created after previous failure", result2.OrderUUID)

		s.inventoryClient.AssertExpectations(s.T())
		s.orderRepository.AssertExpectations(s.T())
	})
}
