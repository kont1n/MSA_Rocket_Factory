package v1

import (
	"context"

	generaredInventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/client/converter"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (c *client) ListParts(ctx context.Context, filter *model.Filter) (*[]model.Part, error) {
	parts, err := c.generatedClient.ListParts(ctx, &generaredInventoryV1.ListPartsRequest{
		Filter: converter.PartsFilterToProto(filter),
	})
	if err != nil {
		return nil, err
	}
	return converter.PartListToModel(parts.Parts), nil
}
