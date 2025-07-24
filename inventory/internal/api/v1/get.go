package v1

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/api/converter"
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	id, err := uuid.Parse(req.PartUuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid part uuid")
	}

	part, err := a.inventoryService.GetPart(ctx, id)
	if err != nil {
		slog.Error("Failed to get part", "id:", id, "error:", err)

		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part not found")
		}

		if errors.Is(err, model.ErrConvertFromRepo) {
			return nil, status.Errorf(codes.Internal, "failed to get part")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	protoPart := converter.ToProtoPart(part)
	return &inventoryV1.GetPartResponse{Part: protoPart}, nil
}
