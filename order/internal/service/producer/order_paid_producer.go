package producer

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	eventsV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/events/v1"
)

type service struct {
	orderPaidProducer kafka.Producer
	metrics           KafkaMetrics
}

// KafkaMetrics интерфейс для метрик Kafka
type KafkaMetrics interface {
	RecordProducerMessage(ctx context.Context, topic string, partition int32, success bool, duration time.Duration)
}

func NewService(orderPaidProducer kafka.Producer, metrics KafkaMetrics) *service {
	return &service{
		orderPaidProducer: orderPaidProducer,
		metrics:           metrics,
	}
}

func (p *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	startTime := time.Now()

	msg := &eventsV1.OrderPaid{
		EventUuid:       event.EventUUID.String(),
		OrderUuid:       event.OrderUUID.String(),
		UserUuid:        event.UserUUID.String(),
		PaymentMethod:   event.PaymentMethod,
		TransactionUuid: event.TransactionUUID.String(),
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "Failed to marshal OrderPaid event", zap.Error(err))
		return err
	}

	err = p.orderPaidProducer.Send(ctx, []byte(event.EventUUID.String()), payload)
	duration := time.Since(startTime)

	// Записываем метрики
	if p.metrics != nil {
		// Для producer'а используем топик "order-paid" и партицию 0 как приблизительное значение
		// В реальной реализации нужно получать реальные значения из SendMessage
		p.metrics.RecordProducerMessage(ctx, "order-paid", 0, err == nil, duration)
	}

	if err != nil {
		logger.Error(ctx, "Failed to send OrderPaid event to Kafka", zap.Error(err))
		return err
	}

	logger.Info(ctx, "OrderPaid event sent successfully",
		zap.String("event_uuid", event.EventUUID.String()),
		zap.String("order_uuid", event.OrderUUID.String()),
		zap.String("payment_method", event.PaymentMethod),
		zap.String("transaction_uuid", event.TransactionUUID.String()))

	return nil
}
