package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	orderV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/openapi/order/v1"
	"github.com/samber/lo"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	orderDraft := model.Order{
		OrderUUID: params.OrderUUID,
	}

	_, err := a.orderService.CancelOrder(ctx, lo.ToPtr(orderDraft))
	if err != nil {
		slog.Error("Cancel order error", "order:", params.OrderUUID, "error:", err)

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
