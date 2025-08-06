package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	orderV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	orderDraft := model.Order{
		UserUUID:  uuid.UUID(req.UserUUID),
		PartUUIDs: req.PartUuids,
	}

	createOrder, err := a.orderService.CreateOrder(ctx, &orderDraft)
	if err != nil {
		slog.Error("Create order error", "order:", req, "error:", err)

		if errors.Is(err, model.ErrPartsSpecified) {
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "parts not specified",
			}, nil
		}

		if errors.Is(err, model.ErrPartsListNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "parts not found",
			}, nil
		}

		if errors.Is(err, model.ErrConvertFromClient) {
			return &orderV1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "can't parse part",
			}, nil
		}

		if errors.Is(err, model.ErrInventoryClient) {
			return &orderV1.InternalServerError{
				Code:    http.StatusFailedDependency,
				Message: "inventory client error",
			}, nil
		}

		return nil, err
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID: orderV1.OrderUUID(createOrder.OrderUUID),
		TotalPrice: orderV1.OptTotalPrice{
			Value: orderV1.TotalPrice(createOrder.TotalPrice),
			Set:   true,
		},
	}, nil
}
