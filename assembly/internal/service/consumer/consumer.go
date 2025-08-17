package consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/kont1n/MSA_Rocket_Factory/assembly/internal/converter/kafka"
	def "github.com/kont1n/MSA_Rocket_Factory/assembly/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	assemblyRecordedConsumer 	kafka.Consumer
	assemblyRecordedDecoder     kafkaConverter.AssemblyRecordedDecoder
}

func NewService(assemblyRecordedConsumer kafka.Consumer, assemblyRecordedDecoder kafkaConverter.AssemblyRecordedDecoder) *service {
	return &service{
		assemblyRecordedConsumer: assemblyRecordedConsumer,
		assemblyRecordedDecoder:  assemblyRecordedDecoder,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting Assembly Consumer service")

	err := s.assemblyRecordedConsumer.Consume(ctx, s.OrderPaidHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}
