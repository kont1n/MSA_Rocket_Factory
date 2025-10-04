package payment

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/tracing"
)

func (s *service) Pay(ctx context.Context, order model.Order) (uuid.UUID, error) {
	// Создаем спан для операции оплаты
	ctx, span := tracing.StartSpan(ctx, "payment.process", trace.WithSpanKind(trace.SpanKindInternal))
	defer tracing.EndSpanWithStatus(span, nil)

	// Добавляем атрибуты к спану
	span.SetAttributes(
		attribute.String("payment.order_uuid", order.OrderUuid.String()),
		attribute.String("payment.user_uuid", order.UserUuid.String()),
		attribute.String("payment.method", order.PaymentMethod),
	)

	span.AddEvent("generating transaction UUID")
	transactionUuid := uuid.New()

	// Добавляем UUID транзакции к спану
	span.SetAttributes(attribute.String("payment.transaction_uuid", transactionUuid.String()))
	span.AddEvent("payment processed successfully", trace.WithAttributes(
		attribute.String("transaction_uuid", transactionUuid.String()),
	))

	logger.Info(ctx, "Payment success",
		zap.String("order_uuid", order.OrderUuid.String()),
		zap.String("transaction_uuid", transactionUuid.String()),
	)

	span.SetStatus(codes.Ok, "payment processed successfully")
	return transactionUuid, nil
}
