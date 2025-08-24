package v1

import (
	"context"
	"fmt"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/client/converter"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	generaredInventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

func (c inventoryClient) ListParts(ctx context.Context, filter *model.Filter) (*[]model.Part, error) {
	parts, err := c.generatedClient.ListParts(ctx, &generaredInventoryV1.ListPartsRequest{
		Filter: converter.ToProtoFilter(filter),
	})
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", err)
	}

	modelParts, err := converter.ToModelPartsList(parts.Parts)
	if err != nil {
		return nil, err
	}

	return modelParts, nil
}
