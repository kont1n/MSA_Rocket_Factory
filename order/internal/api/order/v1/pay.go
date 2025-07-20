package v1

import (
	"context"
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
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "Внутренняя ошибка сервиса - не удалось оплатить заказ",
		}, nil
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: orderV1.TransactionUUID(order.TransactionUUID),
	}, nil
}
