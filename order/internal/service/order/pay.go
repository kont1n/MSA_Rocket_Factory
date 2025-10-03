package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/tracing"
)

func (s service) PayOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	startTime := time.Now()

	// Создаем span для операции оплаты в сервисе
	ctx, span := tracing.StartSpan(ctx, "order.service.pay_order",
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(
			attribute.String("order_uuid", order.OrderUUID.String()),
			attribute.String("payment_method", order.PaymentMethod),
		),
	)
	defer tracing.EndSpanWithStatus(span, nil)

	// Получаем заказ по UUID
	span.AddEvent("fetching order from repository")
	dbOrder, err := s.orderRepository.GetOrder(ctx, order.OrderUUID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.End()
		return nil, fmt.Errorf("service: failed to get order from repository: %w", err)
	}
	span.AddEvent("order fetched successfully")

	// Устанавливаем метод оплаты из запроса
	dbOrder.PaymentMethod = order.PaymentMethod

	// Выполняем запрос к API для оплаты заказа
	span.AddEvent("processing payment via payment client")
	paidOrder, err := s.paymentClient.CreatePayment(ctx, dbOrder)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.End()
		return nil, fmt.Errorf("service: failed to create payment in payment client: %w", err)
	}
	span.AddEvent("payment processed successfully", trace.WithAttributes(
		attribute.String("transaction_uuid", paidOrder.TransactionUUID.String()),
	))

	// Обновляем заказ в хранилище
	span.AddEvent("updating order in repository")
	updatedOrder, err := s.orderRepository.UpdateOrder(ctx, paidOrder)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.End()
		return nil, fmt.Errorf("service: failed to update order in repository: %w", err)
	}

	// Отправляем событие OrderPaid (если producer доступен)
	if s.orderPaidProducer != nil {
		span.AddEvent("producing order paid event")
		event := model.OrderPaidEvent{
			EventUUID:       uuid.New(),
			OrderUUID:       updatedOrder.OrderUUID,
			UserUUID:        updatedOrder.UserUUID,
			PaymentMethod:   updatedOrder.PaymentMethod,
			TransactionUUID: updatedOrder.TransactionUUID,
		}

		err = s.orderPaidProducer.ProduceOrderPaid(ctx, event)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()
			return nil, fmt.Errorf("service: failed to produce OrderPaid event: %w", err)
		}
		span.AddEvent("order paid event produced successfully")
	}

	// Записываем метрики при успешной оплате заказа
	if s.metrics != nil {
		s.metrics.recordOrderPaid(ctx, float64(updatedOrder.TotalPrice), "USD") // По умолчанию USD
		s.metrics.recordOrderDuration(ctx, time.Since(startTime), "pay")
	}

	span.SetStatus(codes.Ok, "order paid successfully")
	return updatedOrder, nil
}
