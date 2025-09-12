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
	*BaseConsumer
	orderPaidDecoder    kafkaConverter.OrderPaidDecoder
	notificationService def.NotificationService
}

func NewOrderPaidService(orderPaidConsumer wrappedKafka.Consumer, orderPaidDecoder kafkaConverter.OrderPaidDecoder, notificationService def.NotificationService, metrics KafkaMetrics) *orderPaidService {
	svc := &orderPaidService{
		orderPaidDecoder:    orderPaidDecoder,
		notificationService: notificationService,
	}

	svc.BaseConsumer = NewBaseConsumer(orderPaidConsumer, svc.OrderPaidHandler, metrics)
	return svc
}

func (s *orderPaidService) OrderPaidHandler(ctx context.Context, msg wrappedKafka.Message) error {
	// Записываем метрики при начале обработки
	if s.metrics != nil {
		s.metrics.RecordConsumerMessage(ctx, "order-paid", 0, "notification-order-paid-group", true)
	}

	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		// Записываем ошибку в метрики
		if s.metrics != nil {
			s.metrics.RecordConsumerMessage(ctx, "order-paid", 0, "notification-order-paid-group", false)
		}
		return err
	}

	err = s.notificationService.NotifyOrderPaid(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send OrderPaid notification", zap.Error(err))
		// Записываем ошибку в метрики
		if s.metrics != nil {
			s.metrics.RecordConsumerMessage(ctx, "order-paid", 0, "notification-order-paid-group", false)
		}
		return err
	}

	logger.Info(ctx, "Successfully processed OrderPaid notification", zap.String("order_uuid", event.OrderUUID.String()))
	return nil
}
