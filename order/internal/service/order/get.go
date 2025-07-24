package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s service) GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	// Получаем заказ по UUID
	order, err := s.orderRepository.GetOrder(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get order from repository: %w", err)
	}

	return order, nil
}
