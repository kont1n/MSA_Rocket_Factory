package order

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s service) GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	// Получаем заказ по UUID
	order, err := s.orderRepository.GetOrder(ctx, id)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "order not found")
	}

	return order, nil
}
