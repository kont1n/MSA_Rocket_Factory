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
	ufoRecordedConsumer kafka.Consumer
	ufoRecordedDecoder  kafkaConverter.UFORecordedDecoder
}

func NewService(ufoRecordedConsumer kafka.Consumer, ufoRecordedDecoder kafkaConverter.UFORecordedDecoder) *service {
	return &service{
		ufoRecordedConsumer: ufoRecordedConsumer,
		ufoRecordedDecoder:  ufoRecordedDecoder,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order ufoRecordedConsumer service")

	err := s.ufoRecordedConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from ufo.recorded topic error", zap.Error(err))
		return err
	}

	return nil
}
