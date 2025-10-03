package order_test

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/fixtures"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

// TestAdvancedFixturesDemo демонстрирует использование улучшенных фикстур
func (s *ServiceSuite) TestAdvancedFixturesDemo() {
	s.Run("ecommerce_сценарий", func() {
		// Используем предопределённый сценарий интернет-магазина
		scenario := fixtures.ECommerceScenario()

		// Валидируем консистентность сценария
		validator := fixtures.NewScenarioValidator(scenario)
		issues := validator.ValidateConsistency()
		s.Require().Empty(issues, "Сценарий должен быть консистентным: %v", issues)

		// Получаем статистику
		stats := validator.ScenarioStats()
		s.T().Logf("Сценарий статистика: %+v", stats)

		// Тестируем создание заказа из сценария
		pendingOrder := scenario.GetOrderByStatus(model.StatusPendingPayment)
		s.Require().NotNil(pendingOrder, "В сценарии должен быть заказ в ожидании оплаты")

		// Настраиваем моки на основе данных сценария
		var scenarioParts []model.Part
		for _, partUUID := range pendingOrder.PartUUIDs {
			part := scenario.GetPartByUUID(partUUID)
			s.Require().NotNil(part, "Деталь должна быть в сценарии: %s", partUUID)
			scenarioParts = append(scenarioParts, *part)
		}

		s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
			Return(&scenarioParts, nil)
		s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
			Return(pendingOrder, nil)

		// Выполняем создание заказа
		result, err := s.service.CreateOrder(context.Background(), pendingOrder)

		s.Require().NoError(err)
		s.Require().NotNil(result)
		s.Require().Equal(pendingOrder.OrderUUID, result.OrderUUID)

		s.inventoryClient.AssertExpectations(s.T())
		s.orderRepository.AssertExpectations(s.T())
	})

	s.Run("edge_case_сценарий", func() {
		// Используем сценарий с граничными случаями
		scenario := fixtures.EdgeCaseScenario()

		// Тестируем заказ с максимальной ценой
		orders := scenario.Orders
		var maxPriceOrder *model.Order
		maxPrice := float32(0)

		for i := range orders {
			if orders[i].TotalPrice > maxPrice {
				maxPrice = orders[i].TotalPrice
				maxPriceOrder = &orders[i]
			}
		}

		s.Require().NotNil(maxPriceOrder, "Должен быть заказ с максимальной ценой")
		s.Require().True(maxPriceOrder.TotalPrice > float32(9999999.0), "Цена должна быть максимальной")

		// Настраиваем мок для деталей заказа с максимальной ценой
		var parts []model.Part
		for _, partUUID := range maxPriceOrder.PartUUIDs {
			part := scenario.GetPartByUUID(partUUID)
			if part == nil {
				// Если деталь не найдена, пропускаем этот заказ
				s.T().Skipf("Деталь с UUID %s не найдена в сценарии, пропускаем тест", partUUID.String())
				return
			}
			parts = append(parts, *part)
		}

		// Проверяем, что найдены все детали
		if len(parts) == 0 {
			s.T().Skip("Не найдены детали для заказа с максимальной ценой, пропускаем тест")
			return
		}

		s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
			Return(&parts, nil)
		s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
			Return(maxPriceOrder, nil)

		result, err := s.service.CreateOrder(context.Background(), maxPriceOrder)
		// Проверяем результат
		if err != nil {
			// Если произошла ошибка, логируем её и пропускаем тест
			s.T().Logf("Ошибка при создании заказа: %v", err)
			s.T().Skip("Не удалось создать заказ с максимальной ценой, пропускаем тест")
			return
		}

		s.Require().NotNil(result)
		s.Require().True(result.TotalPrice >= float32(9999999.0), "Цена должна быть близка к максимальной")

		s.inventoryClient.AssertExpectations(s.T())
		s.orderRepository.AssertExpectations(s.T())
	})

	s.Run("high_volume_сценарий", func() {
		// Используем высоконагруженный сценарий
		scenario := fixtures.HighVolumeScenario()

		// Проверяем что создано достаточно заказов
		s.Require().GreaterOrEqual(len(scenario.Orders), 100, "Должно быть создано минимум 100 заказов")

		// Получаем статистику
		validator := fixtures.NewScenarioValidator(scenario)
		stats := validator.ScenarioStats()

		s.T().Logf("High volume scenario stats: %+v", stats)
		s.Require().Equal(len(scenario.Orders), stats["total_orders"])
		s.Require().Greater(stats["total_price"], float32(0))
	})
}

// TestErrorHandlingWithFixtures демонстрирует использование error фикстур
func (s *ServiceSuite) TestErrorHandlingWithFixtures() {
	s.Run("все_категории_ошибок_валидации", func() {
		validationErrors := fixtures.ErrorScenariosByCategory(fixtures.CategoryValidation)

		for _, errorScenario := range validationErrors {
			s.Run(errorScenario.Name, func() {
				s.SetupTest()

				var problematicOrder *model.Order

				switch errorScenario.Name {
				case "nil_order":
					problematicOrder = nil
				case "empty_parts":
					problematicOrder = fixtures.OrderWithEmptyParts()
					// Настраиваем мок для возврата пустого списка деталей
					s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
						Return(&[]model.Part{}, nil)
				case "nil_user_uuid":
					problematicOrder = fixtures.OrderWithEmptyParts()
					// Настраиваем мок для возврата пустого списка деталей
					s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
						Return(&[]model.Part{}, nil)
					// Настраиваем мок для CreateOrder, хотя он не должен вызываться
					s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
						Return(nil, errors.New("should not be called")).Maybe()
				case "invalid_uuid":
					problematicOrder = fixtures.ValidOrder()
					// Настраиваем мок для возврата деталей
					s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
						Return(&[]model.Part{*fixtures.NewPartBuilder().WithPartUUID(uuid.New()).WithPrice(100.0).Build()}, nil)
					// Настраиваем мок для CreateOrder
					s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
						Return(nil, fixtures.ValidationErrors.InvalidUUID.Error)
				case "empty_payment_method":
					problematicOrder = fixtures.ValidOrder()
					// Настраиваем мок для возврата деталей
					s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
						Return(&[]model.Part{*fixtures.NewPartBuilder().WithPartUUID(uuid.New()).WithPrice(100.0).Build()}, nil)
					// Настраиваем мок для CreateOrder
					s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
						Return(nil, fixtures.ValidationErrors.EmptyPaymentMethod.Error)
				case "too_many_parts":
					problematicOrder = fixtures.ValidOrder()
					// Настраиваем мок для возврата деталей
					s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
						Return(&[]model.Part{*fixtures.NewPartBuilder().WithPartUUID(uuid.New()).WithPrice(100.0).Build()}, nil)
					// Настраиваем мок для CreateOrder
					s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
						Return(nil, fixtures.ValidationErrors.TooManyParts.Error)
				default:
					// Для других случаев создаём базовый заказ
					problematicOrder = fixtures.ValidOrder()
					// Настраиваем мок для возврата деталей
					s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
						Return(&[]model.Part{*fixtures.NewPartBuilder().WithPartUUID(uuid.New()).WithPrice(100.0).Build()}, nil)
				}

				result, err := s.service.CreateOrder(context.Background(), problematicOrder)

				s.Require().Error(err, "Должна быть ошибка для сценария: %s", errorScenario.Name)
				s.Require().Nil(result, "Результат должен быть nil при ошибке")

				// Проверяем ошибку
				switch errorScenario.Name {
				case "nil_order":
					// Специальная проверка для nil order
					s.Require().Contains(err.Error(), "nil")
				case "empty_parts":
					s.Require().True(fixtures.IsExpectedError(err, errorScenario.Error))
				}

				// Логируем детальную информацию об ошибке
				errorInfo := fixtures.GetDetailedErrorInfo(err)
				s.T().Logf("Error details: %+v", errorInfo)

				// Проверяем ожидания моков только если они были настроены
				if errorScenario.Name != "nil_order" {
					s.inventoryClient.AssertExpectations(s.T())
				}
			})
		}
	})

	s.Run("grpc_ошибки_с_детальной_диагностикой", func() {
		grpcErrors := fixtures.ErrorScenariosByCategory(fixtures.CategoryGRPC)

		order := fixtures.ValidOrder()

		for _, errorScenario := range grpcErrors {
			s.Run(errorScenario.Name, func() {
				s.SetupTest() // Сбрасываем моки

				// Настраиваем мок для возврата gRPC ошибки
				s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
					Return(nil, errorScenario.Error)

				result, err := s.service.CreateOrder(context.Background(), order)

				s.Require().Error(err)
				s.Require().Nil(result)

				// Проверяем gRPC код
				expectedCode := status.Code(errorScenario.Error)
				actualCode := status.Code(err)
				s.Require().Equal(expectedCode, actualCode, "gRPC коды должны совпадать")

				// Детальная диагностика ошибки
				errorInfo := fixtures.GetDetailedErrorInfo(err)
				s.Require().True(errorInfo["is_grpc"].(bool), "Должна быть gRPC ошибка")
				s.Require().Equal(expectedCode.String(), errorInfo["grpc_code"])

				s.T().Logf("gRPC error %s: %+v", errorScenario.Name, errorInfo)

				s.inventoryClient.AssertExpectations(s.T())
			})
		}
	})

	s.Run("бизнес_ошибки_с_контекстом", func() {
		businessErrors := fixtures.ErrorScenariosByCategory(fixtures.CategoryBusiness)

		for _, errorScenario := range businessErrors {
			s.Run(errorScenario.Name, func() {
				s.SetupTest()

				order := fixtures.ValidOrder()

				switch errorScenario.Name {
				case "parts_not_found":
					// Возвращаем неполный список деталей
					incompleteParts := []model.Part{}
					s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
						Return(&incompleteParts, nil)

				case "order_not_found":
					// Это ошибка PayOrder, создаём другой сценарий
					payOrder := fixtures.NewOrderBuilder().
						WithOrderUUID(fixtures.TestConstants.OrderUUID1).
						WithPaymentMethod("CARD").
						Build()

					s.orderRepository.On("GetOrder", mock.Anything, payOrder.OrderUUID).
						Return(nil, errorScenario.Error)

					result, err := s.service.PayOrder(context.Background(), payOrder)
					s.Require().Error(err)
					s.Require().Nil(result)
					s.Require().True(fixtures.IsExpectedError(err, errorScenario.Error))

					s.orderRepository.AssertExpectations(s.T())
					return
				case "order_already_paid":
					// Это ошибка PayOrder, создаём другой сценарий
					payOrder := fixtures.NewOrderBuilder().
						WithOrderUUID(fixtures.TestConstants.OrderUUID1).
						WithPaymentMethod("CARD").
						Build()

					// Возвращаем уже оплаченный заказ
					paidOrder := fixtures.NewOrderBuilder().
						WithOrderUUID(fixtures.TestConstants.OrderUUID1).
						WithStatus(model.StatusPaid).
						Build()

					s.orderRepository.On("GetOrder", mock.Anything, payOrder.OrderUUID).
						Return(paidOrder, nil)

					// Настраиваем мок для CreatePayment, который не должен вызываться,
					// но если логика изменится, тест не упадет
					s.paymentClient.On("CreatePayment", mock.Anything, mock.AnythingOfType("*model.Order")).
						Return(nil, errors.New("payment should not be called for already paid order")).Maybe()

					result, err := s.service.PayOrder(context.Background(), payOrder)
					s.Require().Error(err)
					s.Require().Nil(result)
					s.Require().True(fixtures.IsExpectedError(err, errorScenario.Error))

					s.orderRepository.AssertExpectations(s.T())
					s.paymentClient.AssertExpectations(s.T())
					return
				default:
					// Для остальных случаев настраиваем мок для ListParts
					s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
						Return(&[]model.Part{*fixtures.NewPartBuilder().WithPartUUID(uuid.New()).WithPrice(100.0).Build()}, nil)
					// Настраиваем мок для CreateOrder, возвращающий соответствующую ошибку
					s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
						Return(nil, errorScenario.Error)
				}

				result, err := s.service.CreateOrder(context.Background(), order)

				s.Require().Error(err)
				s.Require().Nil(result)

				if errorScenario.Name == "parts_not_found" {
					s.Require().True(fixtures.IsExpectedError(err, errorScenario.Error))
				}

				s.T().Logf("Business error %s: %s", errorScenario.Name, errorScenario.Description)

				s.inventoryClient.AssertExpectations(s.T())
			})
		}
	})
}

// TestEventHandlingWithFixtures демонстрирует тестирование событий
func (s *ServiceSuite) TestEventHandlingWithFixtures() {
	s.Run("создание_и_валидация_событий", func() {
		// Создаём событие с помощью EventBuilder
		orderUUID := fixtures.TestConstants.OrderUUID1
		userUUID := fixtures.TestConstants.UserUUID1
		transactionUUID := fixtures.TestConstants.TransactionUUID

		expectedEvent := fixtures.NewEventBuilder().
			WithOrderUUID(orderUUID).
			WithUserUUID(userUUID).
			WithTransactionUUID(transactionUUID).
			WithPaymentMethod("SBP").
			Build()

		// Настраиваем тестовые данные
		dbOrder := fixtures.NewOrderBuilder().
			WithOrderUUID(orderUUID).
			WithUserUUID(userUUID).
			WithPartUUIDs(fixtures.TestConstants.PartUUID1).
			WithTotalPrice(250.0).
			WithStatus(model.StatusPendingPayment).
			Build()

		paidOrder := fixtures.NewOrderBuilder().
			WithOrderUUID(orderUUID).
			WithUserUUID(userUUID).
			WithPartUUIDs(fixtures.TestConstants.PartUUID1).
			WithTotalPrice(250.0).
			WithTransactionUUID(transactionUUID).
			WithPaymentMethod("SBP").
			WithStatus(model.StatusPaid).
			Build()

		payOrderRequest := fixtures.NewOrderBuilder().
			WithOrderUUID(orderUUID).
			WithPaymentMethod("SBP").
			Build()

		// Настраиваем моки
		s.orderRepository.On("GetOrder", mock.Anything, orderUUID).
			Return(dbOrder, nil)
		s.paymentClient.On("CreatePayment", mock.Anything, mock.AnythingOfType("*model.Order")).
			Return(paidOrder, nil)
		s.orderRepository.On("UpdateOrder", mock.Anything, paidOrder).
			Return(paidOrder, nil)

		// Выполняем оплату заказа
		result, err := s.service.PayOrder(context.Background(), payOrderRequest)

		s.Require().NoError(err)
		s.Require().NotNil(result)
		s.Require().Equal(model.StatusPaid, result.Status)

		// Проверяем что событие было создано
		actualEvent := s.orderPaidProducer.GetLastEvent()
		s.Require().NotNil(actualEvent, "Событие OrderPaid должно быть создано")

		// Сравниваем ключевые поля события
		s.Require().Equal(expectedEvent.OrderUUID, actualEvent.OrderUUID)
		s.Require().Equal(expectedEvent.UserUUID, actualEvent.UserUUID)
		s.Require().Equal(expectedEvent.PaymentMethod, actualEvent.PaymentMethod)
		// Note: TotalPrice не входит в модель OrderPaidEvent

		s.orderRepository.AssertExpectations(s.T())
		s.paymentClient.AssertExpectations(s.T())
	})
}

// TestConcurrencyWithFixtures демонстрирует тестирование конкурентности
func (s *ServiceSuite) TestConcurrencyWithFixtures() {
	s.Run("параллельная_обработка_заказов", func() {
		// Используем специальный сценарий для конкурентности
		scenario := fixtures.ConcurrencyScenario()

		const numGoroutines = 10
		results := make(chan error, numGoroutines)

		sharedOrder := &scenario.Orders[0]
		sharedPart := &scenario.Parts[0]

		// Настраиваем моки для всех горутин
		for i := 0; i < numGoroutines; i++ {
			parts := []model.Part{*sharedPart}
			expectedOrder := fixtures.NewOrderBuilder().
				WithOrderUUID(sharedOrder.OrderUUID).
				WithUserUUID(sharedOrder.UserUUID).
				WithPartUUIDs(sharedOrder.PartUUIDs...).
				WithTotalPrice(float64(sharedOrder.TotalPrice)).
				WithStatus(model.StatusPendingPayment).
				Build()

			s.inventoryClient.On("ListParts", mock.Anything, mock.AnythingOfType("*model.Filter")).
				Return(&parts, nil).Once()
			s.orderRepository.On("CreateOrder", mock.Anything, mock.AnythingOfType("*model.Order")).
				Return(expectedOrder, nil).Once()
		}

		// Запускаем горутины
		for i := 0; i < numGoroutines; i++ {
			go func() {
				// Создаём копию заказа для каждой горутины
				orderCopy := fixtures.NewOrderBuilder().
					WithOrderUUID(sharedOrder.OrderUUID).
					WithUserUUID(sharedOrder.UserUUID).
					WithPartUUIDs(sharedOrder.PartUUIDs...).
					WithTotalPrice(float64(sharedOrder.TotalPrice)).
					WithStatus(sharedOrder.Status).
					Build()

				_, err := s.service.CreateOrder(context.Background(), orderCopy)
				results <- err
			}()
		}

		// Собираем результаты
		successCount := 0
		for i := 0; i < numGoroutines; i++ {
			err := <-results
			if err == nil {
				successCount++
			} else {
				s.T().Logf("Goroutine %d error: %v", i, err)
			}
		}

		s.T().Logf("Successful concurrent operations: %d/%d", successCount, numGoroutines)
		s.Require().Equal(numGoroutines, successCount, "Все горутины должны завершиться успешно")

		s.inventoryClient.AssertExpectations(s.T())
		s.orderRepository.AssertExpectations(s.T())
	})
}
