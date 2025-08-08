package payment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/model"
)

func (s *service) Pay(ctx context.Context, order model.Order) (uuid.UUID, error) {
	transactionUuid := uuid.New()
	slog.Info("Payment success", "transaction_uuid:", transactionUuid)

	return transactionUuid, nil
}
