package consumer

import (
	"context"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
)

// BaseConsumer базовая структура для consumer сервисов
type BaseConsumer struct {
	consumer kafka.Consumer
	handler  MessageHandler
	metrics  KafkaMetrics
}

// KafkaMetrics интерфейс для метрик Kafka
type KafkaMetrics interface {
	RecordConsumerMessage(ctx context.Context, topic string, partition int32, groupID string, success bool)
}

// NewBaseConsumer создает новый базовый consumer
func NewBaseConsumer(consumer kafka.Consumer, handler MessageHandler, metrics KafkaMetrics) *BaseConsumer {
	return &BaseConsumer{
		consumer: consumer,
		handler:  handler,
		metrics:  metrics,
	}
}

// RunConsumer запускает consumer
func (b *BaseConsumer) RunConsumer(ctx context.Context) error {
	return b.consumer.Consume(ctx, b.handler)
}
