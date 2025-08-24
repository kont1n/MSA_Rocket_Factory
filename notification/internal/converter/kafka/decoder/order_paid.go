package decoder

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/model"
	eventsV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/events/v1"
)

type OrderPaidDecoder struct{}

func NewOrderPaidDecoder() *OrderPaidDecoder {
	return &OrderPaidDecoder{}
}

func (d *OrderPaidDecoder) Decode(data []byte) (*model.OrderPaidEvent, error) {
	var protoEvent eventsV1.OrderPaid
	err := proto.Unmarshal(data, &protoEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	eventUUID, err := uuid.Parse(protoEvent.EventUuid)
	if err != nil {
		return nil, fmt.Errorf("failed to parse event UUID: %w", err)
	}

	orderUUID, err := uuid.Parse(protoEvent.OrderUuid)
	if err != nil {
		return nil, fmt.Errorf("failed to parse order UUID: %w", err)
	}

	userUUID, err := uuid.Parse(protoEvent.UserUuid)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user UUID: %w", err)
	}

	transactionUUID, err := uuid.Parse(protoEvent.TransactionUuid)
	if err != nil {
		return nil, fmt.Errorf("failed to parse transaction UUID: %w", err)
	}

	return &model.OrderPaidEvent{
		EventUUID:       eventUUID,
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PaymentMethod:   protoEvent.PaymentMethod,
		TransactionUUID: transactionUUID,
	}, nil
}
