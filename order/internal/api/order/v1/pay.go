package v1

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	orderV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	orderDraft := model.Order{
		OrderUUID:     params.OrderUUID,
		PaymentMethod: string(req.PaymentMethod),
	}

	order, err := a.orderService.PayOrder(ctx, &orderDraft)
	if err != nil {
		logger.Error(ctx, "Pay order error",
			zap.String("order_uuid", params.OrderUUID.String()),
			zap.String("payment_method", string(req.PaymentMethod)),
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

	// Логируем успешную оплату
	logger.Info(ctx, "Order paid successfully",
		zap.String("order_uuid", order.OrderUUID.String()),
		zap.String("payment_method", order.PaymentMethod),
		zap.String("transaction_uuid", order.TransactionUUID.String()),
	)

	return &orderV1.PayOrderResponse{
		TransactionUUID: orderV1.TransactionUUID(order.TransactionUUID),
	}, nil
}
