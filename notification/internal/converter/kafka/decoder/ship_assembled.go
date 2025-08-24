package decoder

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/model"
	eventsV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/events/v1"
)

type ShipAssembledDecoder struct{}

func NewShipAssembledDecoder() *ShipAssembledDecoder {
	return &ShipAssembledDecoder{}
}

func (d *ShipAssembledDecoder) Decode(data []byte) (*model.ShipAssembledEvent, error) {
	var protoEvent eventsV1.ShipAssembled
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

	return &model.ShipAssembledEvent{
		EventUUID: eventUUID,
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		BuildTime: protoEvent.BuildTimeSec,
	}, nil
}
