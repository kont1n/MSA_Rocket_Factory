package order

import (
	"context"
	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s service) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	// Проверяем что детали указаны и заполняем фильтр
	if len(order.PartUUIDs) == 0 {
		return nil, status.Error(codes.FailedPrecondition, "parts not specified")
	}
	uuidFilter := model.Filter{
		PartUUIDs: order.PartUUIDs,
	}

	// Выполняем запрос к API инвентаря для получения деталей заказа
	parts, err := s.inventoryClient.ListParts(ctx, &uuidFilter)
	if err != nil {
		return nil, err
	}
	if len(*parts) != len(order.PartUUIDs) {
		return nil, status.Error(codes.NotFound, "parts not found")
	}

	// Считаем общую стоимость заказа
	totalPrice := 0.0
	for _, part := range *parts {
		totalPrice += part.Price
	}
	order.TotalPrice = totalPrice

	order.OrderUUID = uuid.New()
	order.Status = model.StatusPendingPayment

	// Сохраняем заказ в хранилище
	order, err = s.orderRepository.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}
