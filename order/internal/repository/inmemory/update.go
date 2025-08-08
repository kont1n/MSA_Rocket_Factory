package inmemory

import (
	"context"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/repository/converter"
)

func (r *repository) UpdateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Проверяем, существует ли заказ
	_, exists := r.data[order.OrderUUID.String()]
	if !exists {
		return nil, model.ErrOrderNotFound
	}

	// Конвертируем и сохраняем
	repoOrder := converter.ToRepoOrder(order)
	r.data[order.OrderUUID.String()] = *repoOrder

	return order, nil
}
