package order

import (
	"context"
	"fmt"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s service) CancelOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	// Валидация входных параметров
	if order == nil {
		return nil, fmt.Errorf("order cannot be nil")
	}
	if order.OrderUUID == [16]byte{} {
		return nil, fmt.Errorf("order UUID cannot be nil")
	}

	// Получаем заказ по UUID
	order, err := s.orderRepository.GetOrder(ctx, order.OrderUUID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get order from repository: %w", err)
	}

	// Проверяем статус заказа
	if order.Status == model.StatusPaid {
		return nil, model.ErrPaid
	}
	if order.Status == model.StatusCancelled {
		return nil, model.ErrCancelled
	}

	order.Status = model.StatusCancelled
	// Сохраняем отмену заказа в хранилище
	order, err = s.orderRepository.UpdateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("service: failed to update order in repository: %w", err)
	}

	return order, nil
}
