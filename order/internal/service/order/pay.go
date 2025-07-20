package order

import (
	"context"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s service) PayOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	// Получаем заказ по UUID
	order, err := s.orderRepository.GetOrder(ctx, order.OrderUUID)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос к API для оплаты заказа
	order, err = s.paymentClient.CreatePayment(ctx, order)
	if err != nil {
		return nil, err
	}

	// Обновляем заказ в хранилище
	order, err = s.orderRepository.UpdateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}
