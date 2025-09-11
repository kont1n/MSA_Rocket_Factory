package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/tracing"
)

func (s service) PayOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	startTime := time.Now()

	// Создаем span для операции оплаты в сервисе
	ctx, span := tracing.StartSpan(ctx, "order.service.pay_order",
		trace.WithAttributes(
			attribute.String("order_uuid", order.OrderUUID.String()),
			attribute.String("payment_method", order.PaymentMethod),
		),
	)
	defer span.End()

	// Получаем заказ по UUID
	dbOrder, err := s.orderRepository.GetOrder(ctx, order.OrderUUID)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("service: failed to get order from repository: %w", err)
	}

	// Устанавливаем метод оплаты из запроса
	dbOrder.PaymentMethod = order.PaymentMethod

	// Выполняем запрос к API для оплаты заказа
	paidOrder, err := s.paymentClient.CreatePayment(ctx, dbOrder)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("service: failed to create payment in payment client: %w", err)
	}

	// Обновляем заказ в хранилище
	updatedOrder, err := s.orderRepository.UpdateOrder(ctx, paidOrder)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("service: failed to update order in repository: %w", err)
	}

	// Отправляем событие OrderPaid
	event := model.OrderPaidEvent{
		EventUUID:       uuid.New(),
		OrderUUID:       updatedOrder.OrderUUID,
		UserUUID:        updatedOrder.UserUUID,
		PaymentMethod:   updatedOrder.PaymentMethod,
		TransactionUUID: updatedOrder.TransactionUUID,
	}

	err = s.orderPaidProducer.ProduceOrderPaid(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("service: failed to produce OrderPaid event: %w", err)
	}

	// Записываем метрики при успешной оплате заказа
	if s.metrics != nil {
		s.metrics.recordOrderPaid(ctx, float64(updatedOrder.TotalPrice), "USD") // По умолчанию USD
		s.metrics.recordOrderDuration(ctx, time.Since(startTime), "pay")
	}

	return updatedOrder, nil
}
