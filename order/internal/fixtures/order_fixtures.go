// Package fixtures предоставляет централизованные тестовые данные для order сервиса
package fixtures

import (
	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

// OrderBuilder предоставляет builder pattern для создания тестовых заказов
type OrderBuilder struct {
	order *model.Order
}

// NewOrderBuilder создает новый builder для заказа с дефолтными значениями
func NewOrderBuilder() *OrderBuilder {
	return &OrderBuilder{
		order: &model.Order{
			OrderUUID:       uuid.New(),
			UserUUID:        uuid.New(),
			PartUUIDs:       []uuid.UUID{uuid.New()},
			TotalPrice:      100.0,
			TransactionUUID: uuid.New(),
			PaymentMethod:   "CARD",
			Status:          model.StatusPendingPayment,
		},
	}
}

// WithOrderUUID устанавливает UUID заказа
func (b *OrderBuilder) WithOrderUUID(orderUUID uuid.UUID) *OrderBuilder {
	b.order.OrderUUID = orderUUID
	return b
}

// WithUserUUID устанавливает UUID пользователя
func (b *OrderBuilder) WithUserUUID(userUUID uuid.UUID) *OrderBuilder {
	b.order.UserUUID = userUUID
	return b
}

// WithPartUUIDs устанавливает список деталей
func (b *OrderBuilder) WithPartUUIDs(partUUIDs ...uuid.UUID) *OrderBuilder {
	b.order.PartUUIDs = partUUIDs
	return b
}

// WithTotalPrice устанавливает общую стоимость
func (b *OrderBuilder) WithTotalPrice(price float64) *OrderBuilder {
	b.order.TotalPrice = float32(price)
	return b
}

// WithTransactionUUID устанавливает UUID транзакции
func (b *OrderBuilder) WithTransactionUUID(transactionUUID uuid.UUID) *OrderBuilder {
	b.order.TransactionUUID = transactionUUID
	return b
}

// WithStatus устанавливает статус заказа
func (b *OrderBuilder) WithStatus(status model.OrderStatus) *OrderBuilder {
	b.order.Status = status
	return b
}

// WithPaymentMethod устанавливает способ оплаты
func (b *OrderBuilder) WithPaymentMethod(method string) *OrderBuilder {
	b.order.PaymentMethod = method
	return b
}

// WithCreatedAt removed - model.Order doesn't have CreatedAt field

// EmptyParts очищает список деталей
func (b *OrderBuilder) EmptyParts() *OrderBuilder {
	b.order.PartUUIDs = []uuid.UUID{}
	return b
}

// Build возвращает готовый заказ
func (b *OrderBuilder) Build() *model.Order {
	return b.order
}

// PartBuilder предоставляет builder pattern для создания тестовых деталей
type PartBuilder struct {
	part *model.Part
}

// NewPartBuilder создает новый builder для детали
func NewPartBuilder() *PartBuilder {
	return &PartBuilder{
		part: &model.Part{
			PartUUID: uuid.New(),
			Name:     "Test Part",
			Price:    50.0,
		},
	}
}

// WithPartUUID устанавливает UUID детали
func (b *PartBuilder) WithPartUUID(partUUID uuid.UUID) *PartBuilder {
	b.part.PartUUID = partUUID
	return b
}

// WithPrice устанавливает цену детали
func (b *PartBuilder) WithPrice(price float64) *PartBuilder {
	b.part.Price = price
	return b
}

// WithName устанавливает название детали
func (b *PartBuilder) WithName(name string) *PartBuilder {
	b.part.Name = name
	return b
}

// Build возвращает готовую деталь
func (b *PartBuilder) Build() *model.Part {
	return b.part
}

// Предопределенные фикстуры для часто используемых сценариев

// ValidOrder создает валидный заказ для тестов
func ValidOrder() *model.Order {
	return NewOrderBuilder().Build()
}

// PendingPaymentOrder создает заказ в статусе ожидания оплаты
func PendingPaymentOrder() *model.Order {
	return NewOrderBuilder().
		WithStatus(model.StatusPendingPayment).
		Build()
}

// PaidOrder создает оплаченный заказ
func PaidOrder() *model.Order {
	return NewOrderBuilder().
		WithStatus(model.StatusPaid).
		WithPaymentMethod("CARD").
		Build()
}

// OrderWithEmptyParts создает заказ без деталей для negative тестов
func OrderWithEmptyParts() *model.Order {
	return NewOrderBuilder().
		EmptyParts().
		Build()
}

// TestConstants содержит часто используемые константы для тестов
var TestConstants = struct {
	// Предопределенные UUID для детерминированных тестов
	UserUUID1       uuid.UUID
	UserUUID2       uuid.UUID
	PartUUID1       uuid.UUID
	PartUUID2       uuid.UUID
	PartUUID3       uuid.UUID
	OrderUUID1      uuid.UUID
	OrderUUID2      uuid.UUID
	TransactionUUID uuid.UUID

	// Стандартные значения
	DefaultPrice         float64
	DefaultPaymentMethod string
	DefaultPartPrice     float64
}{
	UserUUID1:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
	UserUUID2:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
	PartUUID1:       uuid.MustParse("660e8400-e29b-41d4-a716-446655440001"),
	PartUUID2:       uuid.MustParse("660e8400-e29b-41d4-a716-446655440002"),
	PartUUID3:       uuid.MustParse("660e8400-e29b-41d4-a716-446655440003"),
	OrderUUID1:      uuid.MustParse("770e8400-e29b-41d4-a716-446655440001"),
	OrderUUID2:      uuid.MustParse("770e8400-e29b-41d4-a716-446655440002"),
	TransactionUUID: uuid.MustParse("880e8400-e29b-41d4-a716-446655440001"),

	DefaultPrice:         100.0,
	DefaultPaymentMethod: "CARD",
	DefaultPartPrice:     50.0,
}

// CreateTestParts создает список тестовых деталей
func CreateTestParts(prices ...float64) []model.Part {
	parts := make([]model.Part, len(prices))
	for i, price := range prices {
		parts[i] = *NewPartBuilder().
			WithPartUUID(uuid.New()).
			WithPrice(price).
			Build()
	}
	return parts
}

// CreateTestPartsWithUUIDs создает список тестовых деталей с предопределенными UUID
func CreateTestPartsWithUUIDs(partUUIDs []uuid.UUID, price float64) []model.Part {
	parts := make([]model.Part, len(partUUIDs))
	for i, partUUID := range partUUIDs {
		parts[i] = *NewPartBuilder().
			WithPartUUID(partUUID).
			WithPrice(price).
			Build()
	}
	return parts
}
