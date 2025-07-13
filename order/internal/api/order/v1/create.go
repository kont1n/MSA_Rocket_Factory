package v1

import (
	"context"
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
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "Внутренняя ошибка сервиса - не удалось получить детали заказа",
		}, nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID: orderV1.OrderUUID(createOrder.OrderUUID),
		TotalPrice: orderV1.OptTotalPrice{
			Value: orderV1.TotalPrice(createOrder.TotalPrice),
			Set:   true,
		},
	}, nil
}
