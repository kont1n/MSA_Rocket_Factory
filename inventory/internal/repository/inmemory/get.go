package inmemory

import (
	"context"

	"github.com/google/uuid"
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
	repoConverter "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/converter"
)

func (r *repository) GetPart(ctx context.Context, uuid uuid.UUID) (*model.Part, error) {
	partUuid := uuid.String()

	r.mu.RLock()
	repoPart, ok := r.data[partUuid]
	r.mu.RUnlock()

	if !ok {
		return nil, model.ErrPartNotFound
	}

	part, err := repoConverter.ToModelPart(repoPart)
	if err != nil {
		return nil, err
	}

	return part, nil
}
