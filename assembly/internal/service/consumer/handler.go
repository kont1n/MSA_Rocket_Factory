package consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

func (s *service) OrderPaidHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.assemblyRecordedDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID.String()),
		zap.String("order_uuid", event.OrderUUID.String()),
	)

	// Вызываем логику сборки корабля
	err = s.assemblyService.Assemble(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to assemble ship", zap.Error(err))
		return err
	}

	return nil
}
