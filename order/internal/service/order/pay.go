package order

import (
	"context"
	"fmt"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s service) PayOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	// Получаем заказ по UUID
	order, err := s.orderRepository.GetOrder(ctx, order.OrderUUID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get order from repository: %w", err)
	}

	// Выполняем запрос к API для оплаты заказа
	order, err = s.paymentClient.CreatePayment(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("service: failed to create payment in payment client: %w", err)
	}

	// Обновляем заказ в хранилище
	order, err = s.orderRepository.UpdateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("service: failed to update order in repository: %w", err)
	}

	return order, nil
}
