package v1

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/api/converter"
	paymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	order, err := converter.PayOrderRqToModel(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	transaction, err := a.paymentService.Pay(ctx, order)
	if err != nil {
		slog.Info("Payment fail", "transaction:", transaction, "err:", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transaction.String(),
	}, nil
}
