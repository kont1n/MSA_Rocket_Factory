package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s service) GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	// Получаем заказ по UUID
	order, err := s.orderRepository.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}
