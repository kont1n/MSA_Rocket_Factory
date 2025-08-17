package consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/kont1n/MSA_Rocket_Factory/notification/internal/converter/kafka"
	def "github.com/kont1n/MSA_Rocket_Factory/notification/internal/service"
	wrappedKafka "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type shipAssembledService struct {
	shipAssembledConsumer wrappedKafka.Consumer
	shipAssembledDecoder  kafkaConverter.ShipAssembledDecoder
	notificationService    def.NotificationService
}

func NewShipAssembledService(shipAssembledConsumer wrappedKafka.Consumer, shipAssembledDecoder kafkaConverter.ShipAssembledDecoder, notificationService def.NotificationService) *shipAssembledService {
	return &shipAssembledService{
		shipAssembledConsumer: shipAssembledConsumer,
		shipAssembledDecoder:  shipAssembledDecoder,
		notificationService:    notificationService,
	}
}

func (s *shipAssembledService) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting ShipAssembled Consumer service")

	err := s.shipAssembledConsumer.Consume(ctx, s.ShipAssembledHandler)
	if err != nil {
		logger.Error(ctx, "Consume from ship.assembled topic error", zap.Error(err))
		return err
	}

	return nil
}

func (s *shipAssembledService) ShipAssembledHandler(ctx context.Context, msg wrappedKafka.Message) error {
	event, err := s.shipAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ShipAssembled", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing ShipAssembled message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID.String()),
		zap.String("order_uuid", event.OrderUUID.String()),
	)

	// Отправляем уведомление
	err = s.notificationService.NotifyShipAssembled(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send ShipAssembled notification", zap.Error(err))
		return err
	}

	return nil
}
