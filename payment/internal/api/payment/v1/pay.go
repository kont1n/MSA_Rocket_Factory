package v1

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/api/converter"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/tracing"
	paymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	// Создаем span для операции оплаты в PaymentService
	ctx, span := tracing.StartSpan(ctx, "payment.pay_order",
		trace.WithAttributes(
			attribute.String("order_uuid", req.OrderUuid),
			attribute.String("payment_method", req.PaymentMethod.String()),
		),
	)
	defer span.End()

	order, err := converter.ToModelOrder(req)
	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	transaction, err := a.paymentService.Pay(ctx, order)
	if err != nil {
		span.RecordError(err)
		logger.Error(ctx, "Payment fail",
			zap.Error(err),
			zap.String("transaction", transaction.String()),
			zap.String("order", order.OrderUuid.String()),
		)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transaction.String(),
	}, nil
}
