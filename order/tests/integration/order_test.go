//go:build integration

package integration

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

var _ = Describe("OrderService", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)
	})

	AfterEach(func() {
		// Чистим таблицу после теста
		err := env.ClearOrdersTable(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку таблицы orders")

		cancel()
	})

	Describe("Репозиторий-слой с PostgreSQL", func() {
		Context("CreateOrder", func() {
			It("должен успешно создавать заказ в PostgreSQL", func() {
				// Создаем тестовый заказ
				testOrder := env.GetTestOrder()

				// Вставляем заказ через тестовое окружение
				orderUUID, err := env.InsertTestOrderWithData(ctx, testOrder)
				Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку заказа")
				Expect(orderUUID).ToNot(BeEmpty(), "ожидали непустой UUID заказа")

				// Проверяем, что заказ действительно создан в базе
				createdOrder, err := env.GetOrderByUUID(ctx, orderUUID)
				Expect(err).ToNot(HaveOccurred(), "ожидали успешное получение заказа")
				Expect(createdOrder).ToNot(BeNil())
				Expect(createdOrder.UserUUID).To(Equal(testOrder.UserUUID))
				Expect(createdOrder.PartUUIDs).To(Equal(testOrder.PartUUIDs))
				Expect(createdOrder.TotalPrice).To(Equal(testOrder.TotalPrice))
				Expect(createdOrder.PaymentMethod).To(Equal(testOrder.PaymentMethod))
				Expect(createdOrder.Status).To(Equal(testOrder.Status))
			})

			It("должен генерировать уникальные UUID для заказов", func() {
				testOrder1 := env.GetTestOrder()
				testOrder2 := env.GetTestOrder()

				orderUUID1, err := env.InsertTestOrderWithData(ctx, testOrder1)
				Expect(err).ToNot(HaveOccurred())

				orderUUID2, err := env.InsertTestOrderWithData(ctx, testOrder2)
				Expect(err).ToNot(HaveOccurred())

				Expect(orderUUID1).ToNot(Equal(orderUUID2), "UUID заказов должны быть уникальными")
			})

			It("должен корректно сохранять массив UUID деталей", func() {
				testOrder := env.GetTestOrder()
				// Добавляем больше деталей для проверки массива
				testOrder.PartUUIDs = []uuid.UUID{
					uuid.New(),
					uuid.New(),
					uuid.New(),
				}

				orderUUID, err := env.InsertTestOrderWithData(ctx, testOrder)
				Expect(err).ToNot(HaveOccurred())

				createdOrder, err := env.GetOrderByUUID(ctx, orderUUID)
				Expect(err).ToNot(HaveOccurred())
				Expect(createdOrder.PartUUIDs).To(HaveLen(3))
				Expect(createdOrder.PartUUIDs).To(Equal(testOrder.PartUUIDs))
			})
		})

		Context("GetOrder", func() {
			var orderUUID string

			BeforeEach(func() {
				// Вставляем тестовый заказ
				var err error
				orderUUID, err = env.InsertTestOrder(ctx)
				Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестового заказа в PostgreSQL")
			})

			It("должен успешно возвращать заказ по UUID", func() {
				order, err := env.GetOrderByUUID(ctx, orderUUID)

				Expect(err).ToNot(HaveOccurred())
				Expect(order).ToNot(BeNil())
				Expect(order.OrderUUID.String()).To(Equal(orderUUID))
				Expect(order.UserUUID).ToNot(Equal(uuid.Nil))
				Expect(order.PartUUIDs).ToNot(BeEmpty())
				Expect(order.TotalPrice).To(BeNumerically(">", 0))
				Expect(order.TransactionUUID).ToNot(Equal(uuid.Nil))
				Expect(order.PaymentMethod).ToNot(BeEmpty())
				Expect(order.Status).ToNot(BeEmpty())
			})

			It("должен возвращать ошибку для несуществующего UUID", func() {
				nonExistentUUID := uuid.New().String()

				order, err := env.GetOrderByUUID(ctx, nonExistentUUID)

				Expect(err).To(HaveOccurred())
				Expect(order).To(BeNil())
				Expect(err.Error()).To(ContainSubstring("заказ не найден"))
			})

			It("должен возвращать ошибку для невалидного UUID", func() {
				invalidUUID := "invalid-uuid-format"

				order, err := env.GetOrderByUUID(ctx, invalidUUID)

				Expect(err).To(HaveOccurred())
				Expect(order).To(BeNil())
			})
		})

		Context("UpdateOrder", func() {
			var orderUUID string
			var originalOrder *model.Order

			BeforeEach(func() {
				// Создаем и вставляем тестовый заказ
				originalOrder = env.GetTestOrder()
				var err error
				orderUUID, err = env.InsertTestOrderWithData(ctx, originalOrder)
				Expect(err).ToNot(HaveOccurred())
			})

			It("должен успешно обновлять статус заказа", func() {
				// Получаем заказ и изменяем его статус
				order, err := env.GetOrderByUUID(ctx, orderUUID)
				Expect(err).ToNot(HaveOccurred())

				// Изменяем статус
				order.Status = model.StatusPaid
				order.PaymentMethod = "bank_transfer"
				order.TransactionUUID = uuid.New()

				// Обновляем заказ через прямое подключение к БД (имитируем работу репозитория)
				err = env.UpdateOrderInDB(ctx, order)
				Expect(err).ToNot(HaveOccurred())

				// Проверяем, что изменения сохранились
				updatedOrder, err := env.GetOrderByUUID(ctx, orderUUID)
				Expect(err).ToNot(HaveOccurred())
				Expect(updatedOrder.Status).To(Equal(model.StatusPaid))
				Expect(updatedOrder.PaymentMethod).To(Equal("bank_transfer"))
				Expect(updatedOrder.TransactionUUID).To(Equal(order.TransactionUUID))
			})

			It("должен сохранять неизменные поля при обновлении", func() {
				order, err := env.GetOrderByUUID(ctx, orderUUID)
				Expect(err).ToNot(HaveOccurred())

				// Сохраняем оригинальные значения
				originalUserUUID := order.UserUUID
				originalPartUUIDs := order.PartUUIDs
				originalTotalPrice := order.TotalPrice

				// Изменяем только статус
				order.Status = model.StatusCancelled

				err = env.UpdateOrderInDB(ctx, order)
				Expect(err).ToNot(HaveOccurred())

				// Проверяем, что остальные поля не изменились
				updatedOrder, err := env.GetOrderByUUID(ctx, orderUUID)
				Expect(err).ToNot(HaveOccurred())
				Expect(updatedOrder.UserUUID).To(Equal(originalUserUUID))
				Expect(updatedOrder.PartUUIDs).To(Equal(originalPartUUIDs))
				Expect(updatedOrder.TotalPrice).To(Equal(originalTotalPrice))
				Expect(updatedOrder.Status).To(Equal(model.StatusCancelled))
			})
		})

		Context("Множественные операции", func() {
			It("должен корректно обрабатывать несколько заказов", func() {
				// Создаем несколько заказов
				orderUUIDs, err := env.InsertMultipleTestOrders(ctx)
				Expect(err).ToNot(HaveOccurred())
				Expect(orderUUIDs).To(HaveLen(3))

				// Проверяем, что все заказы доступны
				for _, orderUUID := range orderUUIDs {
					order, err := env.GetOrderByUUID(ctx, orderUUID)
					Expect(err).ToNot(HaveOccurred())
					Expect(order).ToNot(BeNil())
					Expect(order.OrderUUID.String()).To(Equal(orderUUID))
				}
			})

			It("должен поддерживать различные статусы заказов", func() {
				orderUUIDs, err := env.InsertMultipleTestOrders(ctx)
				Expect(err).ToNot(HaveOccurred())

				// Проверяем, что заказы имеют разные статусы
				var statuses []model.OrderStatus
				for _, orderUUID := range orderUUIDs {
					order, err := env.GetOrderByUUID(ctx, orderUUID)
					Expect(err).ToNot(HaveOccurred())
					statuses = append(statuses, order.Status)
				}

				// Ожидаем разные статусы (pending, paid, cancelled)
				Expect(statuses).To(ContainElement(model.StatusPendingPayment))
				Expect(statuses).To(ContainElement(model.StatusPaid))
				Expect(statuses).To(ContainElement(model.StatusCancelled))
			})
		})

		Context("Граничные случаи и валидация", func() {
			It("должен корректно обрабатывать заказы с одной деталью", func() {
				testOrder := env.GetTestOrder()
				testOrder.PartUUIDs = []uuid.UUID{uuid.New()} // Только одна деталь

				orderUUID, err := env.InsertTestOrderWithData(ctx, testOrder)
				Expect(err).ToNot(HaveOccurred())

				createdOrder, err := env.GetOrderByUUID(ctx, orderUUID)
				Expect(err).ToNot(HaveOccurred())
				Expect(createdOrder.PartUUIDs).To(HaveLen(1))
			})

			It("должен корректно обрабатывать заказы с множественными деталями", func() {
				testOrder := env.GetTestOrder()
				// Создаем заказ с 5 деталями
				testOrder.PartUUIDs = []uuid.UUID{
					uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(),
				}

				orderUUID, err := env.InsertTestOrderWithData(ctx, testOrder)
				Expect(err).ToNot(HaveOccurred())

				createdOrder, err := env.GetOrderByUUID(ctx, orderUUID)
				Expect(err).ToNot(HaveOccurred())
				Expect(createdOrder.PartUUIDs).To(HaveLen(5))
				Expect(createdOrder.PartUUIDs).To(Equal(testOrder.PartUUIDs))
			})

			It("должен корректно обрабатывать различные способы оплаты", func() {
				paymentMethods := []string{"credit_card", "bank_transfer", "cryptocurrency", "cash"}

				for _, method := range paymentMethods {
					testOrder := env.GetTestOrder()
					testOrder.PaymentMethod = method

					orderUUID, err := env.InsertTestOrderWithData(ctx, testOrder)
					Expect(err).ToNot(HaveOccurred())

					createdOrder, err := env.GetOrderByUUID(ctx, orderUUID)
					Expect(err).ToNot(HaveOccurred())
					Expect(createdOrder.PaymentMethod).To(Equal(method))
				}
			})

			It("должен корректно обрабатывать большие суммы заказов", func() {
				testOrder := env.GetTestOrder()
				testOrder.TotalPrice = 999999999.99 // Большая сумма

				orderUUID, err := env.InsertTestOrderWithData(ctx, testOrder)
				Expect(err).ToNot(HaveOccurred())

				createdOrder, err := env.GetOrderByUUID(ctx, orderUUID)
				Expect(err).ToNot(HaveOccurred())
				Expect(createdOrder.TotalPrice).To(BeNumerically("~", testOrder.TotalPrice, 0.01))
			})
		})

		Context("Производительность и нагрузочное тестирование", func() {
			It("должен эффективно обрабатывать множественные операции создания", func() {
				const numOrders = 10
				results := make(chan error, numOrders)

				// Создаем множественные заказы параллельно
				for i := 0; i < numOrders; i++ {
					go func() {
						testOrder := env.GetTestOrder()
						_, err := env.InsertTestOrderWithData(ctx, testOrder)
						results <- err
					}()
				}

				// Проверяем, что все операции выполнились успешно
				for i := 0; i < numOrders; i++ {
					err := <-results
					Expect(err).ToNot(HaveOccurred())
				}
			})

			It("должен эффективно обрабатывать множественные операции чтения", func() {
				// Сначала создаем заказы
				orderUUIDs, err := env.InsertMultipleTestOrders(ctx)
				Expect(err).ToNot(HaveOccurred())

				const numReads = 20
				results := make(chan error, numReads)

				// Выполняем множественные операции чтения параллельно
				for i := 0; i < numReads; i++ {
					orderUUID := orderUUIDs[i%len(orderUUIDs)]
					go func(uuid string) {
						_, err := env.GetOrderByUUID(ctx, uuid)
						results <- err
					}(orderUUID)
				}

				// Проверяем, что все операции выполнились успешно
				for i := 0; i < numReads; i++ {
					err := <-results
					Expect(err).ToNot(HaveOccurred())
				}
			})
		})

		Context("Консистентность данных", func() {
			It("должен возвращать консистентные данные при повторных запросах", func() {
				// Создаем заказ
				orderUUID, err := env.InsertTestOrder(ctx)
				Expect(err).ToNot(HaveOccurred())

				// Выполняем несколько запросов для одного заказа
				var orders []*model.Order
				for i := 0; i < 5; i++ {
					order, err := env.GetOrderByUUID(ctx, orderUUID)
					Expect(err).ToNot(HaveOccurred())
					orders = append(orders, order)
				}

				// Проверяем, что все ответы идентичны
				firstOrder := orders[0]
				for i := 1; i < len(orders); i++ {
					Expect(orders[i].OrderUUID).To(Equal(firstOrder.OrderUUID))
					Expect(orders[i].UserUUID).To(Equal(firstOrder.UserUUID))
					Expect(orders[i].TotalPrice).To(Equal(firstOrder.TotalPrice))
					Expect(orders[i].Status).To(Equal(firstOrder.Status))
					Expect(orders[i].PartUUIDs).To(Equal(firstOrder.PartUUIDs))
				}
			})

			It("должен поддерживать транзакционность операций", func() {
				// Создаем заказ
				testOrder := env.GetTestOrder()
				orderUUID, err := env.InsertTestOrderWithData(ctx, testOrder)
				Expect(err).ToNot(HaveOccurred())

				// Получаем заказ до обновления
				originalOrder, err := env.GetOrderByUUID(ctx, orderUUID)
				Expect(err).ToNot(HaveOccurred())

				// Обновляем заказ
				originalOrder.Status = model.StatusPaid
				err = env.UpdateOrderInDB(ctx, originalOrder)
				Expect(err).ToNot(HaveOccurred())

				// Проверяем, что изменения применились атомарно
				updatedOrder, err := env.GetOrderByUUID(ctx, orderUUID)
				Expect(err).ToNot(HaveOccurred())
				Expect(updatedOrder.Status).To(Equal(model.StatusPaid))
			})
		})
	})

	Describe("Полный сценарий работы с заказами", func() {
		It("должен поддерживать полный жизненный цикл заказа", func() {
			// 1. Создаем новый заказ
			testOrder := env.GetTestOrder()
			testOrder.Status = model.StatusPendingPayment

			orderUUID, err := env.InsertTestOrderWithData(ctx, testOrder)
			Expect(err).ToNot(HaveOccurred())

			// 2. Проверяем, что заказ создан со статусом pending
			createdOrder, err := env.GetOrderByUUID(ctx, orderUUID)
			Expect(err).ToNot(HaveOccurred())
			Expect(createdOrder.Status).To(Equal(model.StatusPendingPayment))

			// 3. Обновляем заказ - переводим в статус paid
			createdOrder.Status = model.StatusPaid
			createdOrder.TransactionUUID = uuid.New()
			err = env.UpdateOrderInDB(ctx, createdOrder)
			Expect(err).ToNot(HaveOccurred())

			// 4. Проверяем, что статус обновился
			paidOrder, err := env.GetOrderByUUID(ctx, orderUUID)
			Expect(err).ToNot(HaveOccurred())
			Expect(paidOrder.Status).To(Equal(model.StatusPaid))
			Expect(paidOrder.TransactionUUID).To(Equal(createdOrder.TransactionUUID))

			// 5. В случае необходимости отменяем заказ
			paidOrder.Status = model.StatusCancelled
			err = env.UpdateOrderInDB(ctx, paidOrder)
			Expect(err).ToNot(HaveOccurred())

			// 6. Проверяем финальный статус
			finalOrder, err := env.GetOrderByUUID(ctx, orderUUID)
			Expect(err).ToNot(HaveOccurred())
			Expect(finalOrder.Status).To(Equal(model.StatusCancelled))
		})

		It("должен поддерживать различные сценарии оплаты", func() {
			paymentScenarios := []struct {
				method string
				status model.OrderStatus
			}{
				{"credit_card", model.StatusPaid},
				{"bank_transfer", model.StatusPaid},
				{"cryptocurrency", model.StatusPaid},
				{"cash", model.StatusPendingPayment},
			}

			for _, scenario := range paymentScenarios {
				testOrder := env.GetTestOrder()
				testOrder.PaymentMethod = scenario.method
				testOrder.Status = scenario.status

				orderUUID, err := env.InsertTestOrderWithData(ctx, testOrder)
				Expect(err).ToNot(HaveOccurred())

				createdOrder, err := env.GetOrderByUUID(ctx, orderUUID)
				Expect(err).ToNot(HaveOccurred())
				Expect(createdOrder.PaymentMethod).To(Equal(scenario.method))
				Expect(createdOrder.Status).To(Equal(scenario.status))
			}
		})
	})
})
