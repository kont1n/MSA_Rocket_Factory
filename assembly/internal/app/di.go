package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/config"
	kafkaConverter "github.com/kont1n/MSA_Rocket_Factory/assembly/internal/converter/kafka"
	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/converter/kafka/decoder"
	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/service"
	assemblyService "github.com/kont1n/MSA_Rocket_Factory/assembly/internal/service/assembly"
	assemblyConsumer "github.com/kont1n/MSA_Rocket_Factory/assembly/internal/service/consumer"
	assemblyProducer "github.com/kont1n/MSA_Rocket_Factory/assembly/internal/service/producer"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/closer"
	wrappedKafka "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka/producer"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	kafkaMiddleware "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	assemblyService         service.AssemblyService
	assemblyProducerService service.ProducerService
	assemblyConsumerService service.ConsumerService

	consumerGroup            sarama.ConsumerGroup
	assemblyRecordedConsumer wrappedKafka.Consumer

	assemblyRecordedDecoder  kafkaConverter.AssemblyRecordedDecoder
	syncProducer             sarama.SyncProducer
	assemblyRecordedProducer wrappedKafka.Producer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

// getKafkaMetrics создает Kafka метрики
func (d *diContainer) getKafkaMetrics(ctx context.Context) *kafkaMiddleware.KafkaMetrics {
	metrics, err := kafkaMiddleware.NewKafkaMetrics()
	if err != nil {
		logger.Error(ctx, "Failed to create Kafka metrics")
		return nil
	}
	return metrics
}

func (d *diContainer) AssemblyService(ctx context.Context) service.AssemblyService {
	if d.assemblyService == nil {
		d.assemblyService = assemblyService.NewService(d.AssemblyProducerService(ctx))
	}

	return d.assemblyService
}

func (d *diContainer) AssemblyProducerService(ctx context.Context) service.ProducerService {
	if d.assemblyProducerService == nil {
		// Получаем Kafka метрики
		kafkaMetrics := d.getKafkaMetrics(ctx)

		d.assemblyProducerService = assemblyProducer.NewService(d.AssemblyRecordedProducer(), kafkaMetrics)
	}

	return d.assemblyProducerService
}

func (d *diContainer) AssemblyConsumerService(ctx context.Context) service.ConsumerService {
	if d.assemblyConsumerService == nil {
		// Получаем Kafka метрики
		kafkaMetrics := d.getKafkaMetrics(ctx)

		d.assemblyConsumerService = assemblyConsumer.NewService(d.AssemblyRecordedConsumer(ctx), d.AssemblyRecordedDecoder(), d.AssemblyService(ctx), kafkaMetrics)
	}

	return d.assemblyConsumerService
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().AssemblyRecordedConsumer.GroupID(),
			config.AppConfig().AssemblyRecordedConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

func (d *diContainer) AssemblyRecordedConsumer(ctx context.Context) wrappedKafka.Consumer {
	if d.assemblyRecordedConsumer == nil {
		// Получаем Kafka метрики
		kafkaMetrics := d.getKafkaMetrics(ctx)

		// Создаем middleware для метрик
		var middlewares []wrappedKafkaConsumer.Middleware
		middlewares = append(middlewares, kafkaMiddleware.Logging(logger.Logger()))

		// Создаем consumer group middleware для отслеживания ребалансировок
		var consumerGroupMiddlewares []wrappedKafka.ConsumerGroupHandlerMiddleware
		if kafkaMetrics != nil {
			rebalancingMiddleware := kafkaMiddleware.RebalancingMiddleware(kafkaMetrics, "assembly-recorded-assembly-group")
			consumerGroupMiddlewares = append(consumerGroupMiddlewares, rebalancingMiddleware)
		}

		d.assemblyRecordedConsumer = wrappedKafkaConsumer.NewConsumerWithConsumerGroupMiddleware(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().AssemblyRecordedConsumer.Topic(),
			},
			logger.Logger(),
			middlewares,
			consumerGroupMiddlewares,
		)
	}

	return d.assemblyRecordedConsumer
}

func (d *diContainer) AssemblyRecordedDecoder() kafkaConverter.AssemblyRecordedDecoder {
	if d.assemblyRecordedDecoder == nil {
		d.assemblyRecordedDecoder = decoder.NewAssemblyRecordedDecoder()
	}

	return d.assemblyRecordedDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().AssemblyRecordedProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) AssemblyRecordedProducer() wrappedKafka.Producer {
	if d.assemblyRecordedProducer == nil {
		d.assemblyRecordedProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().AssemblyRecordedProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.assemblyRecordedProducer
}
