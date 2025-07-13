package order

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s service) CancelOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	// Получаем заказ по UUID
	order, err := s.orderRepository.GetOrder(ctx, order.OrderUUID)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "order not found")
	}

	// Проверяем статус заказа
	if order.Status == model.StatusPaid {
		return nil, status.Error(codes.FailedPrecondition, "order status is paid")
	}
	if order.Status == model.StatusCancelled {
		return nil, status.Error(codes.FailedPrecondition, "order status is cancelled")
	}

	order.Status = model.StatusCancelled
	// Сохраняем отмену заказа в хранилище
	order, err = s.orderRepository.UpdateOrder(ctx, order)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update order")
	}

	return order, nil
}
