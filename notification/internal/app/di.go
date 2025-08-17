package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	telegramClient "github.com/kont1n/MSA_Rocket_Factory/notification/internal/client/telegram"
	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/config"
	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/converter/kafka"
	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/converter/kafka/decoder"
	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/service/consumer"
	notificationService "github.com/kont1n/MSA_Rocket_Factory/notification/internal/service/notification"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/closer"
	wrappedKafka "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka/consumer"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type diContainer struct {
	notificationService        service.NotificationService
	orderPaidConsumer          service.OrderPaidConsumerService
	shipAssembledConsumer      service.ShipAssembledConsumerService
	telegramClient             telegramClient.TelegramClient
	orderPaidConsumerGroup     sarama.ConsumerGroup
	shipAssembledConsumerGroup sarama.ConsumerGroup
	orderPaidKafkaConsumer     wrappedKafka.Consumer
	shipAssembledKafkaConsumer wrappedKafka.Consumer
	orderPaidDecoder           kafka.OrderPaidDecoder
	shipAssembledDecoder       kafka.ShipAssembledDecoder
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) NotificationService(ctx context.Context) service.NotificationService {
	if d.notificationService == nil {
		d.notificationService = notificationService.NewService(d.TelegramClient(ctx))
	}
	return d.notificationService
}

func (d *diContainer) OrderPaidConsumer(ctx context.Context) service.OrderPaidConsumerService {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = consumer.NewOrderPaidService(
			d.OrderPaidKafkaConsumer(),
			d.OrderPaidDecoder(ctx),
			d.NotificationService(ctx),
		)
	}
	return d.orderPaidConsumer
}

func (d *diContainer) ShipAssembledConsumer(ctx context.Context) service.ShipAssembledConsumerService {
	if d.shipAssembledConsumer == nil {
		d.shipAssembledConsumer = consumer.NewShipAssembledService(
			d.ShipAssembledKafkaConsumer(),
			d.ShipAssembledDecoder(ctx),
			d.NotificationService(ctx),
		)
	}
	return d.shipAssembledConsumer
}

func (d *diContainer) TelegramClient(ctx context.Context) telegramClient.TelegramClient {
	if d.telegramClient == nil {
		client, err := telegramClient.NewClient(config.AppConfig().Telegram)
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram client: %s\n", err.Error()))
		}
		d.telegramClient = client
	}
	return d.telegramClient
}

func (d *diContainer) OrderPaidConsumerGroup() sarama.ConsumerGroup {
	if d.orderPaidConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create order paid consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("OrderPaid Kafka consumer group", func(ctx context.Context) error {
			return d.orderPaidConsumerGroup.Close()
		})

		d.orderPaidConsumerGroup = consumerGroup
	}

	return d.orderPaidConsumerGroup
}

func (d *diContainer) ShipAssembledConsumerGroup() sarama.ConsumerGroup {
	if d.shipAssembledConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().ShipAssembledConsumer.GroupID(),
			config.AppConfig().ShipAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create ship assembled consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("ShipAssembled Kafka consumer group", func(ctx context.Context) error {
			return d.shipAssembledConsumerGroup.Close()
		})

		d.shipAssembledConsumerGroup = consumerGroup
	}

	return d.shipAssembledConsumerGroup
}

func (d *diContainer) OrderPaidKafkaConsumer() wrappedKafka.Consumer {
	if d.orderPaidKafkaConsumer == nil {
		d.orderPaidKafkaConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderPaidConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.Logger(),
		)
	}

	return d.orderPaidKafkaConsumer
}

func (d *diContainer) ShipAssembledKafkaConsumer() wrappedKafka.Consumer {
	if d.shipAssembledKafkaConsumer == nil {
		d.shipAssembledKafkaConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ShipAssembledConsumerGroup(),
			[]string{
				config.AppConfig().ShipAssembledConsumer.Topic(),
			},
			logger.Logger(),
		)
	}

	return d.shipAssembledKafkaConsumer
}

func (d *diContainer) OrderPaidDecoder(ctx context.Context) kafka.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) ShipAssembledDecoder(ctx context.Context) kafka.ShipAssembledDecoder {
	if d.shipAssembledDecoder == nil {
		d.shipAssembledDecoder = decoder.NewShipAssembledDecoder()
	}

	return d.shipAssembledDecoder
}
