package part

import (
	"context"
	"fmt"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter *model.Filter) (*[]model.Part, error) {
	parts, err := s.repo.ListParts(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get list of parts from repository: %w", err)
	}

	return parts, nil
}
