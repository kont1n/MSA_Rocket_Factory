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
	assemblyProducerService service.AssemblyProducerService
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

func (d *diContainer) AssemblyService(ctx context.Context) service.AssemblyService {
	if d.assemblyService == nil {
		d.assemblyService = assemblyService.NewService(d.AssemblyProducerService())
	}

	return d.assemblyService
}

func (d *diContainer) AssemblyProducerService() service.AssemblyProducerService {
	if d.assemblyProducerService == nil {
		d.assemblyProducerService = assemblyProducer.NewService(d.AssemblyRecordedProducer())
	}

	return d.assemblyProducerService
}

func (d *diContainer) AssemblyConsumerService() service.ConsumerService {
	if d.assemblyConsumerService == nil {
		d.assemblyConsumerService = assemblyConsumer.NewService(d.AssemblyRecordedConsumer(), d.AssemblyRecordedDecoder())
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

func (d *diContainer) AssemblyRecordedConsumer() wrappedKafka.Consumer {
	if d.assemblyRecordedConsumer == nil {
		d.assemblyRecordedConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().AssemblyRecordedConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
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
