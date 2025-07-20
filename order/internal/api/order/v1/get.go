package v1

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	orderV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	order, err := a.orderService.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		slog.Error("Get order error", "order", params.OrderUUID, "error", err)
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprint("Не удалось найти заказ с таким UUID: ", params.OrderUUID),
		}, nil
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
