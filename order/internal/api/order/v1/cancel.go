package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	orderV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	orderDraft := model.Order{
		OrderUUID: params.OrderUUID,
	}

	_, err := a.orderService.CancelOrder(ctx, lo.ToPtr(orderDraft))
	if err != nil {
		logger.Error(ctx, "Cancel order error",
			zap.String("order_uuid", params.OrderUUID.String()),
			zap.Error(err),
		)

		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "order not found",
			}, nil
		}

		if errors.Is(err, model.ErrConvertFromRepo) {
			return &orderV1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "cannot convert order from repository",
			}, nil
		}

		if errors.Is(err, model.ErrPaid) {
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "order already paid",
			}, nil
		}

		if errors.Is(err, model.ErrCancelled) {
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "order already cancelled",
			}, nil
		}

		return nil, err
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
