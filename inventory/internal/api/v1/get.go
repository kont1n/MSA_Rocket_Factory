package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/api/converter"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	uuid, err := uuid.Parse(req.PartUuid)
	if err != nil {
		return nil, err
	}

	part, err := a.inventoryService.GetPart(ctx, uuid)
	if err != nil {
		return nil, err
	}

	protoPart := converter.ModelToProto(part)
	return &inventoryV1.GetPartResponse{Part: protoPart}, nil
}
