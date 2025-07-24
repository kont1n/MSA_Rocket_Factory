package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s service) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	// Проверяем что детали указаны и заполняем фильтр
	if len(order.PartUUIDs) == 0 {
		return nil, model.ErrPartsSpecified
	}
	uuidFilter := model.Filter{
		PartUUIDs: order.PartUUIDs,
	}

	// Выполняем запрос к API инвентаря для получения деталей заказа
	parts, err := s.inventoryClient.ListParts(ctx, &uuidFilter)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get list parts from inventory client: %w", err)
	}
	if len(*parts) != len(order.PartUUIDs) {
		return nil, model.ErrPartsListNotFound
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
		return nil, fmt.Errorf("service: failed to create order in repository: %w", err)
	}
	return order, nil
}
