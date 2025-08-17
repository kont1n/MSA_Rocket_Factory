package decoder

import (
	"fmt"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/converter"
	"google.golang.org/protobuf/proto"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/model"
	eventsV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewAssemblyRecordedDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsV1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	event, err := converter.ToModelOrder(&pb)
	if err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to convert protobuf to model: %w", err)
	}

	return model.OrderPaidEvent{
		EventUUID:       event.EventUUID,
		OrderUUID:       event.OrderUUID,
		UserUUID:        event.UserUUID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUUID: event.TransactionUUID,
	}, nil
}
