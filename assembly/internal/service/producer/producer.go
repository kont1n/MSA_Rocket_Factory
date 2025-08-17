package producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	eventsV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/events/v1"
)

type service struct {
	assemblyProducer kafka.Producer
}

func NewService(assemblyProducer kafka.Producer) *service {
	return &service{
		assemblyProducer: assemblyProducer,
	}
}

func (p *service) ProduceAssembly(ctx context.Context, event model.ShipAssembledEvent) error {
	msg := &eventsV1.ShipAssembled{
		EventUuid:    event.EventUUID.String(),
		OrderUuid:    event.OrderUUID.String(),
		UserUuid:     event.UserUUID.String(),
		BuildTimeSec: 1,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, model.ErrMarshalToKafkaEvent.Error(), zap.Error(err))
		return err
	}

	err = p.assemblyProducer.Send(ctx, []byte(event.EventUUID.String()), payload)
	if err != nil {
		logger.Error(ctx, model.ErrSendToKafka.Error(), zap.Error(err))
		return err
	}

	return nil
}
