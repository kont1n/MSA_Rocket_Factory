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
	assemblyRecordedConsumer kafka.Consumer
	assemblyRecordedDecoder  kafkaConverter.AssemblyRecordedDecoder
	assemblyService          def.AssemblyService
	metrics                  KafkaMetrics
}

// KafkaMetrics интерфейс для метрик Kafka
type KafkaMetrics interface {
	RecordConsumerMessage(ctx context.Context, topic string, partition int32, groupID string, success bool)
}

func NewService(assemblyRecordedConsumer kafka.Consumer, assemblyRecordedDecoder kafkaConverter.AssemblyRecordedDecoder, assemblyService def.AssemblyService, metrics KafkaMetrics) *service {
	return &service{
		assemblyRecordedConsumer: assemblyRecordedConsumer,
		assemblyRecordedDecoder:  assemblyRecordedDecoder,
		assemblyService:          assemblyService,
		metrics:                  metrics,
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

func (s *service) OrderPaidHandler(ctx context.Context, msg kafka.Message) error {
	// Записываем метрики при начале обработки
	if s.metrics != nil {
		s.metrics.RecordConsumerMessage(ctx, "order-paid", 0, "assembly-recorded-assembly-group", true)
	}

	// Декодируем сообщение
	event, err := s.assemblyRecordedDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid event", zap.Error(err))
		// Записываем ошибку в метрики
		if s.metrics != nil {
			s.metrics.RecordConsumerMessage(ctx, "order-paid", 0, "assembly-recorded-assembly-group", false)
		}
		return err
	}

	// Обрабатываем событие
	err = s.assemblyService.Assemble(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to process OrderPaid event", zap.Error(err))
		// Записываем ошибку в метрики
		if s.metrics != nil {
			s.metrics.RecordConsumerMessage(ctx, "order-paid", 0, "assembly-recorded-assembly-group", false)
		}
		return err
	}

	logger.Info(ctx, "Successfully processed OrderPaid event", zap.String("order_uuid", event.OrderUUID.String()))
	return nil
}
