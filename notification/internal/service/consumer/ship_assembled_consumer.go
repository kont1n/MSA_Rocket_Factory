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
	*BaseConsumer
	shipAssembledDecoder kafkaConverter.ShipAssembledDecoder
	notificationService  def.NotificationService
}

func NewShipAssembledService(shipAssembledConsumer wrappedKafka.Consumer, shipAssembledDecoder kafkaConverter.ShipAssembledDecoder, notificationService def.NotificationService, metrics KafkaMetrics) *shipAssembledService {
	svc := &shipAssembledService{
		shipAssembledDecoder: shipAssembledDecoder,
		notificationService:  notificationService,
	}

	svc.BaseConsumer = NewBaseConsumer(shipAssembledConsumer, svc.ShipAssembledHandler, metrics)
	return svc
}

func (s *shipAssembledService) ShipAssembledHandler(ctx context.Context, msg wrappedKafka.Message) error {
	// Записываем метрики при начале обработки
	if s.metrics != nil {
		s.metrics.RecordConsumerMessage(ctx, "ship-assembled", 0, "notification-ship-assembled-group", true)
	}

	event, err := s.shipAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ShipAssembled", zap.Error(err))
		// Записываем ошибку в метрики
		if s.metrics != nil {
			s.metrics.RecordConsumerMessage(ctx, "ship-assembled", 0, "notification-ship-assembled-group", false)
		}
		return err
	}

	err = s.notificationService.NotifyShipAssembled(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send ShipAssembled notification", zap.Error(err))
		// Записываем ошибку в метрики
		if s.metrics != nil {
			s.metrics.RecordConsumerMessage(ctx, "ship-assembled", 0, "notification-ship-assembled-group", false)
		}
		return err
	}

	logger.Info(ctx, "Successfully processed ShipAssembled notification", zap.String("order_uuid", event.OrderUUID.String()))
	return nil
}
