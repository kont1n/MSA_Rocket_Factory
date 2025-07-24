package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
)

type InventoryService interface {
	GetPart(ctx context.Context, uuid uuid.UUID) (*model.Part, error)
	ListParts(ctx context.Context, filter *model.Filter) (*[]model.Part, error)
}
