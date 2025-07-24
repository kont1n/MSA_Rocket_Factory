package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	orderV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	order, err := a.orderService.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		slog.Error("Get order error", "order", params.OrderUUID, "error", err)

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

		return nil, err
	}

	return &orderV1.OrderDto{
		OrderUUID: order.OrderUUID,
		UserUUID:  order.UserUUID,
		PartUuids: order.PartUUIDs,
		TotalPrice: orderV1.OptFloat32{
			Value: float32(order.TotalPrice),
			Set:   true,
		},
		TransactionUUID: orderV1.OptUUID{
			Value: order.TransactionUUID,
			Set:   true,
		},
		PaymentMethod: orderV1.OptPaymentMethod{
			Value: orderV1.PaymentMethod(order.PaymentMethod),
			Set:   true,
		},
		Status: orderV1.OrderStatus(order.Status),
	}, nil
}
