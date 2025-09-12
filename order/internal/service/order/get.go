package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/tracing"
)

func (s service) GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	// Создаем спан для операции получения заказа
	ctx, span := tracing.StartSpan(ctx, "order.get")
	defer span.End()

	// Добавляем атрибуты к спану
	span.SetAttributes(attribute.String("order.uuid", id.String()))

	// Получаем заказ по UUID
	order, err := s.orderRepository.GetOrder(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("service: failed to get order from repository: %w", err)
	}

	// Добавляем информацию о найденном заказе к спану
	span.SetAttributes(
		attribute.String("order.status", string(order.Status)),
		attribute.Float64("order.total_price", float64(order.TotalPrice)),
	)

	return order, nil
}
