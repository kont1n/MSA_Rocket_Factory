// Package fixtures предоставляет расширенные фикстуры для сложных сценариев тестирования
package fixtures

import (
	"time"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

// ScenarioBuilder предоставляет builder для создания сложных тестовых сценариев
type ScenarioBuilder struct {
	orders []model.Order
	parts  []model.Part
	events []model.OrderPaidEvent
}

// NewScenarioBuilder создаёт новый builder для сценариев
func NewScenarioBuilder() *ScenarioBuilder {
	return &ScenarioBuilder{
		orders: make([]model.Order, 0),
		parts:  make([]model.Part, 0),
		events: make([]model.OrderPaidEvent, 0),
	}
}

// WithOrders добавляет заказы к сценарию
func (b *ScenarioBuilder) WithOrders(orders ...*model.Order) *ScenarioBuilder {
	for _, order := range orders {
		if order != nil {
			b.orders = append(b.orders, *order)
		}
	}
	return b
}

// WithParts добавляет детали к сценарию
func (b *ScenarioBuilder) WithParts(parts ...*model.Part) *ScenarioBuilder {
	for _, part := range parts {
		if part != nil {
			b.parts = append(b.parts, *part)
		}
	}
	return b
}

// WithEvents добавляет события к сценарию
func (b *ScenarioBuilder) WithEvents(events ...*model.OrderPaidEvent) *ScenarioBuilder {
	for _, event := range events {
		if event != nil {
			b.events = append(b.events, *event)
		}
	}
	return b
}

// Build возвращает построенный сценарий
func (b *ScenarioBuilder) Build() *TestScenario {
	return &TestScenario{
		Orders: b.orders,
		Parts:  b.parts,
		Events: b.events,
	}
}

// TestScenario представляет комплексный тестовый сценарий
type TestScenario struct {
	Orders []model.Order
	Parts  []model.Part
	Events []model.OrderPaidEvent
}

// GetOrderByStatus возвращает первый заказ с указанным статусом
func (ts *TestScenario) GetOrderByStatus(status model.OrderStatus) *model.Order {
	for i := range ts.Orders {
		if ts.Orders[i].Status == status {
			return &ts.Orders[i]
		}
	}
	return nil
}

// GetPartByUUID возвращает деталь по UUID
func (ts *TestScenario) GetPartByUUID(partUUID uuid.UUID) *model.Part {
	for i := range ts.Parts {
		if ts.Parts[i].PartUUID == partUUID {
			return &ts.Parts[i]
		}
	}
	return nil
}

// EventBuilder предоставляет builder для создания событий OrderPaid
type EventBuilder struct {
	event *model.OrderPaidEvent
}

// NewEventBuilder создаёт новый builder для событий
func NewEventBuilder() *EventBuilder {
	return &EventBuilder{
		event: &model.OrderPaidEvent{
			OrderUUID:       uuid.New(),
			UserUUID:        uuid.New(),
			TransactionUUID: uuid.New(),
			PaymentMethod:   "CARD",
		},
	}
}

// WithOrderUUID устанавливает UUID заказа в событии
func (b *EventBuilder) WithOrderUUID(orderUUID uuid.UUID) *EventBuilder {
	b.event.OrderUUID = orderUUID
	return b
}

// WithUserUUID устанавливает UUID пользователя в событии
func (b *EventBuilder) WithUserUUID(userUUID uuid.UUID) *EventBuilder {
	b.event.UserUUID = userUUID
	return b
}

// WithTransactionUUID устанавливает UUID транзакции в событии
func (b *EventBuilder) WithTransactionUUID(transactionUUID uuid.UUID) *EventBuilder {
	b.event.TransactionUUID = transactionUUID
	return b
}

// WithPaymentMethod устанавливает метод оплаты в событии
func (b *EventBuilder) WithPaymentMethod(method string) *EventBuilder {
	b.event.PaymentMethod = method
	return b
}

// WithTotalPrice удалён - поле TotalPrice отсутствует в model.OrderPaidEvent

// Build возвращает готовое событие
func (b *EventBuilder) Build() *model.OrderPaidEvent {
	return b.event
}

// Предопределённые сложные сценарии

// ECommerceScenario создаёт типичный сценарий интернет-магазина
func ECommerceScenario() *TestScenario {
	userUUID := TestConstants.UserUUID1

	// Создаём детали для разных категорий
	enginePart := NewPartBuilder().
		WithPartUUID(TestConstants.PartUUID1).
		WithName("Rocket Engine V1").
		WithPrice(5000.0).
		Build()

	fuelTankPart := NewPartBuilder().
		WithPartUUID(TestConstants.PartUUID2).
		WithName("Fuel Tank 500L").
		WithPrice(1500.0).
		Build()

	navigationPart := NewPartBuilder().
		WithPartUUID(TestConstants.PartUUID3).
		WithName("Navigation System").
		WithPrice(2000.0).
		Build()

	// Создаём заказы в разных статусах
	pendingOrder := NewOrderBuilder().
		WithOrderUUID(TestConstants.OrderUUID1).
		WithUserUUID(userUUID).
		WithPartUUIDs(TestConstants.PartUUID1, TestConstants.PartUUID2).
		WithTotalPrice(6500.0).
		WithStatus(model.StatusPendingPayment).
		WithPaymentMethod("CARD").
		Build()

	paidOrder := NewOrderBuilder().
		WithOrderUUID(TestConstants.OrderUUID2).
		WithUserUUID(userUUID).
		WithPartUUIDs(TestConstants.PartUUID3).
		WithTotalPrice(2000.0).
		WithStatus(model.StatusPaid).
		WithPaymentMethod("SBP").
		WithTransactionUUID(TestConstants.TransactionUUID).
		Build()

	// Создаём соответствующее событие
	paidEvent := NewEventBuilder().
		WithOrderUUID(TestConstants.OrderUUID2).
		WithUserUUID(userUUID).
		WithTransactionUUID(TestConstants.TransactionUUID).
		WithPaymentMethod("SBP").
		Build()

	return NewScenarioBuilder().
		WithOrders(pendingOrder, paidOrder).
		WithParts(enginePart, fuelTankPart, navigationPart).
		WithEvents(paidEvent).
		Build()
}

// HighVolumeScenario создаёт сценарий с большим количеством заказов
func HighVolumeScenario() *TestScenario {
	builder := NewScenarioBuilder()

	// Создаём 100 заказов
	for i := 0; i < 100; i++ {
		orderUUID := uuid.New()
		userUUID := uuid.New()
		partUUID := uuid.New()

		// Создаём деталь
		part := NewPartBuilder().
			WithPartUUID(partUUID).
			WithName("Part " + orderUUID.String()[:8]).
			WithPrice(float64((i + 1) * 10)).
			Build()

		// Создаём заказ
		status := model.StatusPendingPayment
		if i%3 == 0 {
			status = model.StatusPaid
		} else if i%5 == 0 {
			status = model.StatusCancelled
		}

		order := NewOrderBuilder().
			WithOrderUUID(orderUUID).
			WithUserUUID(userUUID).
			WithPartUUIDs(partUUID).
			WithTotalPrice(float64((i + 1) * 10)).
			WithStatus(status).
			Build()

		builder = builder.WithOrders(order).WithParts(part)

		// Добавляем события для оплаченных заказов
		if status == model.StatusPaid {
			event := NewEventBuilder().
				WithOrderUUID(orderUUID).
				WithUserUUID(userUUID).
				WithTransactionUUID(uuid.New()).
				WithPaymentMethod("CARD").
				Build()
			builder = builder.WithEvents(event)
		}
	}

	return builder.Build()
}

// EdgeCaseScenario создаёт сценарий с граничными случаями
func EdgeCaseScenario() *TestScenario {
	// Заказ с максимальной ценой
	maxPriceOrder := NewOrderBuilder().
		WithOrderUUID(uuid.New()).
		WithUserUUID(TestConstants.UserUUID1).
		WithPartUUIDs(TestConstants.PartUUID1).
		WithTotalPrice(9999999.99).
		WithStatus(model.StatusPendingPayment).
		Build()

	// Заказ с минимальной ценой
	minPriceOrder := NewOrderBuilder().
		WithOrderUUID(uuid.New()).
		WithUserUUID(TestConstants.UserUUID2).
		WithPartUUIDs(TestConstants.PartUUID2).
		WithTotalPrice(0.01).
		WithStatus(model.StatusPaid).
		Build()

	// Заказ с максимальным количеством деталей
	maxPartsUUIDs := make([]uuid.UUID, 1000)
	for i := range maxPartsUUIDs {
		maxPartsUUIDs[i] = uuid.New()
	}

	maxPartsOrder := NewOrderBuilder().
		WithOrderUUID(uuid.New()).
		WithUserUUID(TestConstants.UserUUID1).
		WithPartUUIDs(maxPartsUUIDs...).
		WithTotalPrice(50000.0).
		WithStatus(model.StatusPendingPayment).
		Build()

	// Соответствующие детали
	maxPricePart := NewPartBuilder().
		WithPartUUID(TestConstants.PartUUID1).
		WithPrice(9999999.99).
		WithName("Ultra Expensive Part").
		Build()

	minPricePart := NewPartBuilder().
		WithPartUUID(TestConstants.PartUUID2).
		WithPrice(0.01).
		WithName("Ultra Cheap Part").
		Build()

	return NewScenarioBuilder().
		WithOrders(maxPriceOrder, minPriceOrder, maxPartsOrder).
		WithParts(maxPricePart, minPricePart).
		Build()
}

// ConcurrencyScenario создаёт сценарий для тестирования конкурентности
func ConcurrencyScenario() *TestScenario {
	builder := NewScenarioBuilder()

	// Создаём один заказ, который будут пытаться обработать параллельно
	sharedOrderUUID := uuid.New()
	sharedUserUUID := uuid.New()
	sharedPartUUID := uuid.New()

	sharedPart := NewPartBuilder().
		WithPartUUID(sharedPartUUID).
		WithName("Shared Part").
		WithPrice(100.0).
		Build()

	sharedOrder := NewOrderBuilder().
		WithOrderUUID(sharedOrderUUID).
		WithUserUUID(sharedUserUUID).
		WithPartUUIDs(sharedPartUUID).
		WithTotalPrice(100.0).
		WithStatus(model.StatusPendingPayment).
		Build()

	return builder.
		WithOrders(sharedOrder).
		WithParts(sharedPart).
		Build()
}

// ErrorRecoveryScenario создаёт сценарий для тестирования восстановления после ошибок
func ErrorRecoveryScenario() *TestScenario {
	// Заказ который изначально не может быть создан из-за временной ошибки
	problematicOrder := NewOrderBuilder().
		WithOrderUUID(uuid.New()).
		WithUserUUID(TestConstants.UserUUID1).
		WithPartUUIDs(TestConstants.PartUUID1).
		WithTotalPrice(100.0).
		WithStatus(model.StatusPendingPayment).
		Build()

	// Деталь которая становится доступной после "исправления" проблемы
	recoverablePart := NewPartBuilder().
		WithPartUUID(TestConstants.PartUUID1).
		WithName("Recoverable Part").
		WithPrice(100.0).
		Build()

	return NewScenarioBuilder().
		WithOrders(problematicOrder).
		WithParts(recoverablePart).
		Build()
}

// Утилиты для работы со сценариями

// ScenarioValidator предоставляет валидацию сценариев
type ScenarioValidator struct {
	scenario *TestScenario
}

// NewScenarioValidator создаёт новый валидатор
func NewScenarioValidator(scenario *TestScenario) *ScenarioValidator {
	return &ScenarioValidator{scenario: scenario}
}

// ValidateConsistency проверяет консистентность данных в сценарии
func (v *ScenarioValidator) ValidateConsistency() []string {
	var issues []string

	// Проверяем что для каждого заказа есть соответствующие детали
	for _, order := range v.scenario.Orders {
		for _, partUUID := range order.PartUUIDs {
			if v.scenario.GetPartByUUID(partUUID) == nil {
				issues = append(issues, "Part not found for order "+order.OrderUUID.String()+": "+partUUID.String())
			}
		}
	}

	// Проверяем что для каждого события есть соответствующий заказ
	for _, event := range v.scenario.Events {
		found := false
		for _, order := range v.scenario.Orders {
			if order.OrderUUID == event.OrderUUID {
				found = true
				// Проверяем что статус заказа соответствует событию
				if order.Status != model.StatusPaid {
					issues = append(issues, "Order status inconsistent with paid event: "+order.OrderUUID.String())
				}
				break
			}
		}
		if !found {
			issues = append(issues, "Order not found for paid event: "+event.OrderUUID.String())
		}
	}

	return issues
}

// ScenarioStats возвращает статистику по сценарию
func (v *ScenarioValidator) ScenarioStats() map[string]interface{} {
	statusCounts := make(map[model.OrderStatus]int)
	totalPrice := float32(0)
	totalParts := 0

	for _, order := range v.scenario.Orders {
		statusCounts[order.Status]++
		totalPrice += order.TotalPrice
		totalParts += len(order.PartUUIDs)
	}

	return map[string]interface{}{
		"total_orders":        len(v.scenario.Orders),
		"total_parts":         len(v.scenario.Parts),
		"total_events":        len(v.scenario.Events),
		"status_counts":       statusCounts,
		"total_price":         totalPrice,
		"avg_parts_per_order": float64(totalParts) / float64(len(v.scenario.Orders)),
	}
}

// TimeBasedBuilder создаёт заказы с временными метками
type TimeBasedBuilder struct {
	baseTime time.Time
	interval time.Duration
	counter  int
}

// NewTimeBasedBuilder создаёт builder с временными метками
func NewTimeBasedBuilder(baseTime time.Time, interval time.Duration) *TimeBasedBuilder {
	return &TimeBasedBuilder{
		baseTime: baseTime,
		interval: interval,
		counter:  0,
	}
}

// NextOrder создаёт следующий заказ с инкрементальной временной меткой
func (tb *TimeBasedBuilder) NextOrder() *OrderBuilder {
	// Вычисляем время для потенциального использования в будущем
	_ = tb.baseTime.Add(time.Duration(tb.counter) * tb.interval)
	tb.counter++

	return NewOrderBuilder().
		WithOrderUUID(uuid.New()).
		WithUserUUID(TestConstants.UserUUID1)
	// Note: model.Order не имеет полей времени, но здесь показан паттерн для будущего расширения
}

// Reset сбрасывает счётчик
func (tb *TimeBasedBuilder) Reset() {
	tb.counter = 0
}
