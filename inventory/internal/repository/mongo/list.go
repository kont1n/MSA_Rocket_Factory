package mongo

import (
	"context"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
)

func (r *repository) ListParts(ctx context.Context, filter *model.Filter) (*[]model.Part, error) {
	// TODO implement me
	return nil, nil
}
