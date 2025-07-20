package part

import (
	"context"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter *model.Filter) (*[]model.Part, error) {
	parts, err := s.repo.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}

	return parts, nil
}
