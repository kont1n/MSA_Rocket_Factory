package producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type service struct {
	ufoRecordedProducer kafka.Producer
}

func NewService(ufoRecordedProducer kafka.Producer) *service {
	return &service{
		ufoRecordedProducer: ufoRecordedProducer,
	}
}

func (p *service) ProduceAssemblyRecorded(ctx context.Context, event model.AssemblyRecordedEvent) error {
	var observedAt *timestamppb.Timestamp
	if event.ObservedAt != nil {
		observedAt = timestamppb.New(*event.ObservedAt)
	}

	msg := &eventsV1.UFORecorded{
		ObservedAt:  observedAt,
		Location:    event.Location,
		Description: event.Description,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal UFORecorded", zap.Error(err))
		return err
	}

	err = p.ufoRecordedProducer.Send(ctx, []byte(event.UUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish UFORecorded", zap.Error(err))
		return err
	}

	return nil
}
