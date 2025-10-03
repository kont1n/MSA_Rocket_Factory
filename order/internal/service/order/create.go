package order

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/tracing"
)

func (s service) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	// Создаем спан для операции создания заказа
	ctx, span := tracing.StartSpan(ctx, "order.create", trace.WithSpanKind(trace.SpanKindInternal))
	defer func() {
		// Паника будет обработана выше, здесь просто устанавливаем статус
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("panic: %v", r)
			}
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()
			panic(r)
		}
	}()

	startTime := time.Now()

	// Валидация входных параметров
	span.AddEvent("validating order parameters")
	if order == nil {
		err := fmt.Errorf("order cannot be nil")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.End()
		return nil, err
	}
	if order.UserUUID == [16]byte{} {
		err := fmt.Errorf("user UUID cannot be nil")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.End()
		return nil, err
	}

	// Проверяем что детали указаны и заполняем фильтр
	if len(order.PartUUIDs) == 0 {
		span.RecordError(model.ErrPartsSpecified)
		span.SetStatus(codes.Error, model.ErrPartsSpecified.Error())
		span.End()
		return nil, model.ErrPartsSpecified
	}
	// Проверяем лимит количества деталей
	if len(order.PartUUIDs) > 1000 {
		err := fmt.Errorf("too many parts in order: maximum 1000 allowed")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.End()
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
	span.AddEvent("fetching parts from inventory")
	parts, err := s.inventoryClient.ListParts(ctx, &uuidFilter)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.End()
		return nil, fmt.Errorf("service: failed to get list parts from inventory client: %w", err)
	}
	if len(*parts) != len(order.PartUUIDs) {
		span.RecordError(model.ErrPartsListNotFound)
		span.SetStatus(codes.Error, model.ErrPartsListNotFound.Error())
		span.End()
		return nil, model.ErrPartsListNotFound
	}
	span.AddEvent("parts fetched successfully", trace.WithAttributes(attribute.Int("parts_count", len(*parts))))

	// Считаем общую стоимость заказа
	span.AddEvent("calculating total price")
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
	span.AddEvent("saving order to repository")
	order, err = s.orderRepository.CreateOrder(ctx, order)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.End()
		return nil, fmt.Errorf("service: failed to create order in repository: %w", err)
	}

	// Добавляем UUID созданного заказа к спану
	span.SetAttributes(attribute.String("order.uuid", order.OrderUUID.String()))
	span.AddEvent("order created successfully", trace.WithAttributes(attribute.String("order_uuid", order.OrderUUID.String())))

	// Записываем метрики при успешном создании заказа
	if s.metrics != nil {
		s.metrics.recordOrderCreated(ctx, totalPrice, "USD") // По умолчанию USD
		s.metrics.recordOrderDuration(ctx, time.Since(startTime), "create")
	}

	span.SetStatus(codes.Ok, "order created successfully")
	span.End()
	return order, nil
}
