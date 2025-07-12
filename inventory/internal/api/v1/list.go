package v1

import (
	"context"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/api/converter"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filter := converter.ProtoToModel(req)

	parts, err := a.inventoryService.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}

	protoParts := make([]*inventoryV1.Part, 0, len(*parts))
	for _, part := range *parts {
		protoPart := converter.ModelToProto(&part)
		protoParts = append(protoParts, protoPart)
	}

	return &inventoryV1.ListPartsResponse{Parts: protoParts}, nil
}
