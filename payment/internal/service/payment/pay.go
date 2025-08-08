package payment

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

func (s *service) Pay(ctx context.Context, order model.Order) (uuid.UUID, error) {
	transactionUuid := uuid.New()
	logger.Info(ctx, "Payment success",
		zap.String("transaction_uuid", transactionUuid.String()),
	)

	return transactionUuid, nil
}
