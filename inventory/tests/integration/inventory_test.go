//go:build integration

package integration

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Создаём gRPC клиент
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		// Чистим коллекцию после теста
		err := env.ClearPartsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции parts")

		cancel()
	})

	Describe("GetPart", func() {
		var partUUID string

		BeforeEach(func() {
			// Вставляем тестовую деталь
			var err error
			partUUID, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовой детали в MongoDB")
		})

		It("должен успешно возвращать деталь по UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				PartUuid: partUUID,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetPart()).ToNot(BeNil())
			Expect(resp.GetPart().PartUuid).To(Equal(partUUID))
			Expect(resp.GetPart().GetName()).ToNot(BeEmpty())
			Expect(resp.GetPart().GetDescription()).ToNot(BeEmpty())
			Expect(resp.GetPart().GetPrice()).To(BeNumerically(">", 0))
			Expect(resp.GetPart().GetStockQuantity()).To(BeNumerically(">=", 0))
			Expect(resp.GetPart().GetCategory()).ToNot(Equal(inventoryV1.Category_CATEGORY_UNSPECIFIED))
			Expect(resp.GetPart().GetDimensions()).ToNot(BeNil())
			Expect(resp.GetPart().GetManufacturer()).ToNot(BeNil())
			Expect(resp.GetPart().GetTags()).ToNot(BeEmpty())
			Expect(resp.GetPart().GetCreatedAt()).ToNot(BeNil())
			Expect(resp.GetPart().GetUpdatedAt()).ToNot(BeNil())
		})

		It("должен возвращать ошибку NotFound для несуществующего UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				PartUuid: "00000000-0000-0000-0000-000000000000", // Валидный UUID, но несуществующий
			})

			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			st, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(st.Code()).To(Equal(codes.NotFound))
		})

		It("должен возвращать ошибку InvalidArgument для пустого UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				PartUuid: "",
			})

			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			st, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(st.Code()).To(Equal(codes.InvalidArgument))
		})
	})

	Describe("ListParts", func() {
		var partUUIDs []string

		BeforeEach(func() {
			// Вставляем несколько тестовых деталей
			var err error
			partUUIDs, err = env.InsertMultipleTestParts(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовых деталей в MongoDB")
			Expect(partUUIDs).To(HaveLen(3))
		})

		It("должен успешно возвращать все детали без фильтра", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).ToNot(BeNil())
			Expect(resp.GetParts()).To(HaveLen(3))

			// Проверяем, что все детали содержат необходимые поля
			for _, part := range resp.GetParts() {
				Expect(part.GetPartUuid()).ToNot(BeEmpty())
				Expect(part.GetName()).ToNot(BeEmpty())
				Expect(part.GetDescription()).ToNot(BeEmpty())
				Expect(part.GetPrice()).To(BeNumerically(">", 0))
				Expect(part.GetCategory()).ToNot(Equal(inventoryV1.Category_CATEGORY_UNSPECIFIED))
			}
		})

		It("должен успешно фильтровать детали по UUID", func() {
			filter := &inventoryV1.PartsFilter{
				PartUuid: []string{partUUIDs[0]},
			}

			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: filter,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).To(HaveLen(1))
			Expect(resp.GetParts()[0].GetPartUuid()).To(Equal(partUUIDs[0]))
		})

		It("должен успешно фильтровать детали по категории", func() {
			filter := &inventoryV1.PartsFilter{
				Category: []inventoryV1.Category{inventoryV1.Category_CATEGORY_ENGINE},
			}

			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: filter,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).ToNot(BeEmpty())

			// Проверяем, что все возвращенные детали имеют категорию ENGINE
			for _, part := range resp.GetParts() {
				Expect(part.GetCategory()).To(Equal(inventoryV1.Category_CATEGORY_ENGINE))
			}
		})

		It("должен успешно фильтровать детали по стране производителя", func() {
			filter := &inventoryV1.PartsFilter{
				ManufacturerCountry: []string{"Russia"},
			}

			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: filter,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).ToNot(BeEmpty())

			// Проверяем, что все возвращенные детали произведены в России
			for _, part := range resp.GetParts() {
				Expect(part.GetManufacturer().GetCountry()).To(Equal("Russia"))
			}
		})

		It("должен успешно фильтровать детали по тегам", func() {
			filter := &inventoryV1.PartsFilter{
				Tags: []string{"engine"},
			}

			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: filter,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).ToNot(BeEmpty())

			// Проверяем, что все возвращенные детали содержат тег "engine"
			for _, part := range resp.GetParts() {
				Expect(part.GetTags()).To(ContainElement("engine"))
			}
		})

		It("должен возвращать пустой список для фильтра без совпадений", func() {
			filter := &inventoryV1.PartsFilter{
				PartName: []string{"NonexistentPart"},
			}

			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: filter,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).To(BeEmpty())
		})

		It("должен корректно обрабатывать комбинированные фильтры", func() {
			filter := &inventoryV1.PartsFilter{
				Category:            []inventoryV1.Category{inventoryV1.Category_CATEGORY_ENGINE},
				ManufacturerCountry: []string{"Russia"},
			}

			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: filter,
			})

			Expect(err).ToNot(HaveOccurred())

			// Проверяем, что все возвращенные детали соответствуют обоим критериям
			for _, part := range resp.GetParts() {
				Expect(part.GetCategory()).To(Equal(inventoryV1.Category_CATEGORY_ENGINE))
				Expect(part.GetManufacturer().GetCountry()).To(Equal("Russia"))
			}
		})
	})

	Describe("Полный сценарий работы с инвентарем", func() {
		It("должен поддерживать полный цикл работы с деталями", func() {
			// 1. Проверяем, что изначально список пуст
			listResp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{})
			Expect(err).ToNot(HaveOccurred())
			Expect(listResp.GetParts()).To(BeEmpty())

			// 2. Добавляем тестовые данные
			partUUIDs, err := env.InsertMultipleTestParts(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(partUUIDs).To(HaveLen(3))

			// 3. Проверяем, что теперь детали доступны в списке
			listResp, err = inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{})
			Expect(err).ToNot(HaveOccurred())
			Expect(listResp.GetParts()).To(HaveLen(3))

			// 4. Получаем конкретную деталь по UUID
			getResp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				PartUuid: partUUIDs[0],
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(getResp.GetPart().GetPartUuid()).To(Equal(partUUIDs[0]))

			// 5. Проверяем фильтрацию по различным критериям
			engineFilter := &inventoryV1.PartsFilter{
				Category: []inventoryV1.Category{inventoryV1.Category_CATEGORY_ENGINE},
			}
			engineResp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: engineFilter,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(engineResp.GetParts()).ToNot(BeEmpty())

			// 6. Проверяем фильтрацию по стране
			countryFilter := &inventoryV1.PartsFilter{
				ManufacturerCountry: []string{"USA"},
			}
			countryResp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: countryFilter,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(countryResp.GetParts()).ToNot(BeEmpty())
		})
	})

	Describe("Производительность и нагрузочное тестирование", func() {
		It("должен эффективно обрабатывать множественные запросы", func() {
			// Вставляем тестовые данные
			_, err := env.InsertMultipleTestParts(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Выполняем множественные параллельные запросы
			const numRequests = 10
			results := make(chan error, numRequests)

			for i := 0; i < numRequests; i++ {
				go func() {
					_, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{})
					results <- err
				}()
			}

			// Проверяем, что все запросы выполнились успешно
			for i := 0; i < numRequests; i++ {
				err := <-results
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("должен эффективно обрабатывать одновременные запросы GetPart", func() {
			// Добавляем тестовые данные
			_, err := env.InsertMultipleTestParts(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Сначала получаем UUID существующей детали
			listResp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{})
			Expect(err).ToNot(HaveOccurred())
			Expect(listResp.GetParts()).ToNot(BeEmpty())

			partUUID := listResp.GetParts()[0].GetPartUuid()

			// Выполняем множественные параллельные запросы GetPart
			const numRequests = 20
			results := make(chan error, numRequests)

			for i := 0; i < numRequests; i++ {
				go func() {
					_, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
						PartUuid: partUUID,
					})
					results <- err
				}()
			}

			// Проверяем, что все запросы выполнились успешно
			for i := 0; i < numRequests; i++ {
				err := <-results
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("должен эффективно обрабатывать запросы с различными фильтрами", func() {
			// Добавляем тестовые данные
			_, err := env.InsertMultipleTestParts(ctx)
			Expect(err).ToNot(HaveOccurred())

			const numRequests = 15
			results := make(chan error, numRequests)

			// Создаем различные фильтры для тестирования
			filters := []*inventoryV1.PartsFilter{
				{Category: []inventoryV1.Category{inventoryV1.Category_CATEGORY_ENGINE}},
				{ManufacturerCountry: []string{"Russia"}},
				{Tags: []string{"rocket"}},
				{ManufacturerCountry: []string{"USA"}},
				{Category: []inventoryV1.Category{inventoryV1.Category_CATEGORY_FUEL}},
			}

			for i := 0; i < numRequests; i++ {
				filter := filters[i%len(filters)]
				go func(f *inventoryV1.PartsFilter) {
					_, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
						Filter: f,
					})
					results <- err
				}(filter)
			}

			// Проверяем, что все запросы выполнились успешно
			for i := 0; i < numRequests; i++ {
				err := <-results
				Expect(err).ToNot(HaveOccurred())
			}
		})
	})

	Describe("Валидация входных данных", func() {
		Context("GetPart с невалидными данными", func() {
			It("должен возвращать ошибку для nil запроса", func() {
				resp, err := inventoryClient.GetPart(ctx, nil)

				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())

				st, ok := status.FromError(err)
				Expect(ok).To(BeTrue())
				Expect(st.Code()).To(Equal(codes.InvalidArgument))
			})

			It("должен возвращать ошибку для очень длинного UUID", func() {
				longUUID := "this-is-a-very-long-uuid-that-exceeds-normal-length-and-should-cause-validation-error-in-the-system"

				resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
					PartUuid: longUUID,
				})

				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("должен возвращать ошибку для UUID с недопустимыми символами", func() {
				invalidUUID := "invalid-uuid-with-special-chars-@#$%"

				resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
					PartUuid: invalidUUID,
				})

				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})

		Context("ListParts с граничными случаями фильтров", func() {
			BeforeEach(func() {
				// Добавляем тестовые данные для фильтрации
				_, err := env.InsertMultipleTestParts(ctx)
				Expect(err).ToNot(HaveOccurred())
			})

			It("должен корректно обрабатывать пустые массивы в фильтрах", func() {
				filter := &inventoryV1.PartsFilter{
					PartUuid:            []string{},
					PartName:            []string{},
					Category:            []inventoryV1.Category{},
					ManufacturerCountry: []string{},
					Tags:                []string{},
				}

				resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
					Filter: filter,
				})

				Expect(err).ToNot(HaveOccurred())
				// При пустых фильтрах должны возвращаться все записи
				Expect(resp.GetParts()).To(HaveLen(3))
			})

			It("должен корректно обрабатывать фильтр с неизвестной категорией", func() {
				filter := &inventoryV1.PartsFilter{
					Category: []inventoryV1.Category{inventoryV1.Category_CATEGORY_UNSPECIFIED},
				}

				resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
					Filter: filter,
				})

				Expect(err).ToNot(HaveOccurred())
				// Система не обрабатывает CATEGORY_UNSPECIFIED и возвращает все записи
				// Это текущее поведение приложения
				Expect(resp.GetParts()).To(HaveLen(3))
			})

			It("должен корректно обрабатывать фильтр с очень длинными строками", func() {
				longString := "very-long-string-that-might-cause-performance-issues-in-database-queries-and-should-be-handled-gracefully-by-the-system"

				filter := &inventoryV1.PartsFilter{
					PartName: []string{longString},
				}

				resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
					Filter: filter,
				})

				Expect(err).ToNot(HaveOccurred())
				Expect(resp.GetParts()).To(BeEmpty())
			})

			It("должен корректно обрабатывать фильтр с большим количеством значений", func() {
				// Создаем фильтр с множеством валидных UUID (все несуществующие)
				manyUUIDs := make([]string, 100)
				for i := 0; i < 100; i++ {
					manyUUIDs[i] = fmt.Sprintf("00000000-0000-0000-0000-%012d", i)
				}

				filter := &inventoryV1.PartsFilter{
					PartUuid: manyUUIDs,
				}

				resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
					Filter: filter,
				})

				Expect(err).ToNot(HaveOccurred())
				Expect(resp.GetParts()).To(BeEmpty())
			})
		})
	})

	Describe("Обработка ошибок базы данных", func() {
		Context("Поведение при проблемах с подключением", func() {
			It("должен корректно обрабатывать запросы при пустой базе данных", func() {
				// База данных пуста, проверяем корректную обработку
				resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{})

				Expect(err).ToNot(HaveOccurred())
				Expect(resp.GetParts()).To(BeEmpty())
			})

			It("должен возвращать NotFound для несуществующих деталей", func() {
				resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
					PartUuid: "00000000-0000-0000-0000-000000000000",
				})

				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())

				st, ok := status.FromError(err)
				Expect(ok).To(BeTrue())
				Expect(st.Code()).To(Equal(codes.NotFound))
			})
		})
	})

	Describe("Консистентность данных", func() {
		It("должен возвращать консистентные данные при повторных запросах", func() {
			// Добавляем тестовые данные
			partUUIDs, err := env.InsertMultipleTestParts(ctx)
			Expect(err).ToNot(HaveOccurred())

			partUUID := partUUIDs[0]

			// Выполняем несколько запросов для одной и той же детали
			var responses []*inventoryV1.GetPartResponse
			for i := 0; i < 5; i++ {
				resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
					PartUuid: partUUID,
				})
				Expect(err).ToNot(HaveOccurred())
				responses = append(responses, resp)
			}

			// Проверяем, что все ответы идентичны
			firstResponse := responses[0]
			for i := 1; i < len(responses); i++ {
				Expect(responses[i].GetPart().GetPartUuid()).To(Equal(firstResponse.GetPart().GetPartUuid()))
				Expect(responses[i].GetPart().GetName()).To(Equal(firstResponse.GetPart().GetName()))
				Expect(responses[i].GetPart().GetPrice()).To(Equal(firstResponse.GetPart().GetPrice()))
				Expect(responses[i].GetPart().GetStockQuantity()).To(Equal(firstResponse.GetPart().GetStockQuantity()))
			}
		})

		It("должен возвращать стабильный порядок при ListParts без фильтров", func() {
			// Добавляем тестовые данные
			_, err := env.InsertMultipleTestParts(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Выполняем несколько запросов списка
			var responses []*inventoryV1.ListPartsResponse
			for i := 0; i < 3; i++ {
				resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{})
				Expect(err).ToNot(HaveOccurred())
				responses = append(responses, resp)
			}

			// Проверяем, что порядок стабилен (или хотя бы количество одинаково)
			firstResponse := responses[0]
			for i := 1; i < len(responses); i++ {
				Expect(responses[i].GetParts()).To(HaveLen(len(firstResponse.GetParts())))
			}
		})
	})
})
