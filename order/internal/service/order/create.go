package order

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/tracing"
)

func (s service) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	// Создаем спан для операции создания заказа
	ctx, span := tracing.StartSpan(ctx, "order.create")
	defer span.End()

	startTime := time.Now()

	// Валидация входных параметров
	if order == nil {
		span.RecordError(fmt.Errorf("order cannot be nil"))
		return nil, fmt.Errorf("order cannot be nil")
	}
	if order.UserUUID == [16]byte{} {
		span.RecordError(fmt.Errorf("user UUID cannot be nil"))
		return nil, fmt.Errorf("user UUID cannot be nil")
	}

	// Проверяем что детали указаны и заполняем фильтр
	if len(order.PartUUIDs) == 0 {
		span.RecordError(model.ErrPartsSpecified)
		return nil, model.ErrPartsSpecified
	}
	// Проверяем лимит количества деталей
	if len(order.PartUUIDs) > 1000 {
		err := fmt.Errorf("too many parts in order: maximum 1000 allowed")
		span.RecordError(err)
		return nil, err
	}

	// Добавляем атрибуты к спану
	span.SetAttributes(
		attribute.String("order.user_uuid", order.UserUUID.String()),
		attribute.Int("order.parts_count", len(order.PartUUIDs)),
	)

	uuidFilter := model.Filter{
		PartUUIDs: order.PartUUIDs,
	}

	// Выполняем запрос к API инвентаря для получения деталей заказа
	parts, err := s.inventoryClient.ListParts(ctx, &uuidFilter)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("service: failed to get list parts from inventory client: %w", err)
	}
	if len(*parts) != len(order.PartUUIDs) {
		span.RecordError(model.ErrPartsListNotFound)
		return nil, model.ErrPartsListNotFound
	}

	// Считаем общую стоимость заказа
	totalPrice := 0.0
	for _, part := range *parts {
		totalPrice += part.Price
	}
	order.TotalPrice = float32(totalPrice)

	// Добавляем информацию о стоимости к спану
	span.SetAttributes(
		attribute.Float64("order.total_price", totalPrice),
		attribute.String("order.status", string(model.StatusPendingPayment)),
	)

	order.Status = model.StatusPendingPayment

	// Сохраняем заказ в хранилище
	order, err = s.orderRepository.CreateOrder(ctx, order)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("service: failed to create order in repository: %w", err)
	}

	// Добавляем UUID созданного заказа к спану
	span.SetAttributes(attribute.String("order.uuid", order.OrderUUID.String()))

	// Записываем метрики при успешном создании заказа
	if s.metrics != nil {
		s.metrics.recordOrderCreated(ctx, totalPrice, "USD") // По умолчанию USD
		s.metrics.recordOrderDuration(ctx, time.Since(startTime), "create")
	}

	return order, nil
}
