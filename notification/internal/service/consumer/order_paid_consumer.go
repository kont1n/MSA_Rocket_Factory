package consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/kont1n/MSA_Rocket_Factory/notification/internal/converter/kafka"
	def "github.com/kont1n/MSA_Rocket_Factory/notification/internal/service"
	wrappedKafka "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type orderPaidService struct {
	orderPaidConsumer   wrappedKafka.Consumer
	orderPaidDecoder    kafkaConverter.OrderPaidDecoder
	notificationService def.NotificationService
}

func NewOrderPaidService(orderPaidConsumer wrappedKafka.Consumer, orderPaidDecoder kafkaConverter.OrderPaidDecoder, notificationService def.NotificationService) *orderPaidService {
	return &orderPaidService{
		orderPaidConsumer:   orderPaidConsumer,
		orderPaidDecoder:    orderPaidDecoder,
		notificationService: notificationService,
	}
}

func (s *orderPaidService) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting OrderPaid Consumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderPaidHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}

func (s *orderPaidService) OrderPaidHandler(ctx context.Context, msg wrappedKafka.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing OrderPaid message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID.String()),
		zap.String("order_uuid", event.OrderUUID.String()),
	)

	// Отправляем уведомление
	err = s.notificationService.NotifyOrderPaid(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send OrderPaid notification", zap.Error(err))
		return err
	}

	return nil
}
