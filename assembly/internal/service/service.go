package service

import (
	"context"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/model"
)

type AssemblyService interface {
	Assemble()
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type AssemblyProducerService interface {
	ProduceAssemblyRecorded(ctx context.Context, event model.AssemblyRecordedEvent) error
}
