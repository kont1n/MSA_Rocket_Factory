package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	orderV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	orderDraft := model.Order{
		OrderUUID:     params.OrderUUID,
		PaymentMethod: string(req.PaymentMethod),
	}

	order, err := a.orderService.PayOrder(ctx, &orderDraft)
	if err != nil {
		slog.Error("Pay order error", "order", params.OrderUUID, "error", err)

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

		if errors.Is(err, model.ErrConvertFromClient) {
			return &orderV1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "can't parse transaction",
			}, nil
		}

		if errors.Is(err, model.ErrPaymentClient) {
			return &orderV1.InternalServerError{
				Code:    http.StatusFailedDependency,
				Message: "payment client error",
			}, nil
		}

		return nil, err
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: orderV1.TransactionUUID(order.TransactionUUID),
	}, nil
}
