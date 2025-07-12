package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuid uuid.UUID) (*model.Part, error) {
	part, err := s.repo.GetPart(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return part, nil
}