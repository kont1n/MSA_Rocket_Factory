package inmemory

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/repository/converter"
)

func (r *repository) GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	r.mu.RLock()
	repoOrder := r.data[id.String()]
	r.mu.RUnlock()

	if lo.ToPtr(repoOrder) == nil {
		return nil, model.ErrPartNotFound
	}

	order, err := converter.RepoToModel(lo.ToPtr(repoOrder))
	if err != nil {
		return nil, err
	}

	return order, nil
}
