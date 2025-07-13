package v1

import (
	"context"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	orderV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	orderDraft := model.Order{
		OrderUUID: params.OrderUUID,
	}

	_, err := a.orderService.CancelOrder(ctx, &orderDraft)
	if err != nil {
		return nil, err
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
