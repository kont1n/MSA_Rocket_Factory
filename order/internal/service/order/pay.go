package order

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s service) PayOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	// Получаем заказ по UUID
	order, err := s.orderRepository.GetOrder(ctx, order.OrderUUID)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "order not found")
	}

	// Выполняем запрос к API для оплаты заказа
	order, err = s.paymentClient.CreatePayment(ctx, order)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "failed to pay for the order ")
	}

	// Обновляем заказ в хранилище
	order, err = s.orderRepository.UpdateOrder(ctx, order)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "failed to update order ")
	}

	return order, nil
}
