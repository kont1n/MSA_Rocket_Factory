package mongo

import (
	"context"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
)

func (r *repository) GetPart(ctx context.Context, uuid uuid.UUID) (*model.Part, error) {
	//TODO implement me
	return nil, nil
}
