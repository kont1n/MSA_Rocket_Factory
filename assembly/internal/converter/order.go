package converter

import (
	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/model"
	eventsV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/events/v1"
)

func ToModelOrder (order *eventsV1.OrderPaid) (*model.OrderPaidEvent, error) {
	eventId, err := uuid.Parse(order.EventUuid)
	if err != nil {
		return nil, model.ErrConvertFromKafkaEvent
	}

	orderId, err := uuid.Parse(order.OrderUuid)
	if err != nil {
		return nil, model.ErrConvertFromKafkaEvent
	}

	userId, err := uuid.Parse(order.UserUuid)
	if err != nil {
		return nil, model.ErrConvertFromKafkaEvent
	}

	transactionId, err := uuid.Parse(order.TransactionUuid)
	if err != nil {
		return nil, model.ErrConvertFromKafkaEvent
	}

	return &model.OrderPaidEvent{
		EventUUID:       eventId,
		OrderUUID:       orderId,
		UserUUID:        userId,
		PaymentMethod:   order.PaymentMethod,
		TransactionUUID: transactionId,
	}, nil
}
