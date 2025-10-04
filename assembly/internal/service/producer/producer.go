package producer

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	eventsV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/events/v1"
)

type service struct {
	assemblyProducer kafka.Producer
	metrics          KafkaMetrics
}

// KafkaMetrics интерфейс для метрик Kafka
type KafkaMetrics interface {
	RecordProducerMessage(ctx context.Context, topic string, partition int32, success bool, duration time.Duration)
}

func NewService(assemblyProducer kafka.Producer, metrics KafkaMetrics) *service {
	return &service{
		assemblyProducer: assemblyProducer,
		metrics:          metrics,
	}
}

func (p *service) ProduceAssembly(ctx context.Context, event model.ShipAssembledEvent) error {
	startTime := time.Now()

	msg := &eventsV1.ShipAssembled{
		EventUuid:    event.EventUUID.String(),
		OrderUuid:    event.OrderUUID.String(),
		UserUuid:     event.UserUUID.String(),
		BuildTimeSec: event.BuildTime,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, model.ErrMarshalToKafkaEvent.Error(), zap.Error(err))
		return err
	}

	err = p.assemblyProducer.Send(ctx, []byte(event.EventUUID.String()), payload)
	duration := time.Since(startTime)

	// Записываем метрики
	if p.metrics != nil {
		p.metrics.RecordProducerMessage(ctx, "ship-assembled", 0, err == nil, duration)
	}

	if err != nil {
		logger.Error(ctx, model.ErrSendToKafka.Error(), zap.Error(err))
		return err
	}

	return nil
}
