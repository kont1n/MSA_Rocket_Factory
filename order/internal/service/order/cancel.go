package order

import (
	"context"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s service) CancelOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	// Получаем заказ по UUID
	order, err := s.orderRepository.GetOrder(ctx, order.OrderUUID)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return order, nil
}
