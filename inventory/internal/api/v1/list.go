package v1

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/api/converter"
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filter := converter.ToModelPart(req)

	parts, err := a.inventoryService.ListParts(ctx, filter)
	if err != nil {
		slog.Error("Failed to get list part", "filter:", filter, "error:", err)

		if errors.Is(err, model.ErrConvertFromRepo) {
			return nil, status.Errorf(codes.Internal, "failed to get list parts")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	protoParts := make([]*inventoryV1.Part, 0, len(*parts))
	for _, part := range *parts {
		protoPart := converter.ToProtoPart(&part)
		protoParts = append(protoParts, protoPart)
	}

	return &inventoryV1.ListPartsResponse{Parts: protoParts}, nil
}
