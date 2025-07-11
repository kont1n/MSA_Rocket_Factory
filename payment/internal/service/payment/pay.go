package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/model"
)

func (s *service) Pay(ctx context.Context, order model.Order) (uuid.UUID, error) {
	transaction_uuid := uuid.New()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s\n", transaction_uuid)

	return transaction_uuid, nil
}
