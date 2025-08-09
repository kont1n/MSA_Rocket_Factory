package inmemory

import (
	"context"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/repository/converter"
)

func (r *repository) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	// Генерируем новый UUID для заказа
	order.OrderUUID = uuid.New()

	// Конвертируем в repo модель ПОСЛЕ установки UUID
	repoOrder := converter.ToRepoOrder(order)

	r.mu.Lock()
	r.data[order.OrderUUID.String()] = *repoOrder
	r.mu.Unlock()

	return order, nil
}
