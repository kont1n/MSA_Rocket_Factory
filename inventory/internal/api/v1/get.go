package v1

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/api/converter"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	uuid, err := uuid.Parse(req.PartUuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid part uuid")
	}

	part, err := a.inventoryService.GetPart(ctx, uuid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get part")
	}
	if part == nil {
		return nil, status.Errorf(codes.NotFound, "part not found")
	}

	protoPart := converter.ModelToProto(part)
	return &inventoryV1.GetPartResponse{Part: protoPart}, nil
}
