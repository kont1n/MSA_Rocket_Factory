package app

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	iamClient "github.com/kont1n/MSA_Rocket_Factory/notification/internal/client/iam"
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
	iamV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/iam/v1"
)

type diContainer struct {
	mu                         sync.RWMutex
	telegramClientLocked       bool
	notificationService        service.NotificationService
	orderPaidConsumer          service.OrderPaidConsumerService
	shipAssembledConsumer      service.ShipAssembledConsumerService
	telegramClient             telegramClient.TelegramClient
	iamClient                  iamClient.Client
	iamGRPCConn                *grpc.ClientConn
	iamGRPCClient              iamV1.UserServiceClient
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
	d.mu.RLock()
	if d.notificationService != nil {
		defer d.mu.RUnlock()
		return d.notificationService
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.notificationService != nil {
		return d.notificationService
	}

	// Сначала создаем IAM клиент, затем Telegram клиент
	iamClientInstance := d.IAMClient(ctx)

	d.telegramClientLocked = true
	telegramClientInstance := d.telegramClientUnsafe(ctx)
	d.telegramClientLocked = false

	d.notificationService = notificationService.NewService(ctx, telegramClientInstance, iamClientInstance)
	return d.notificationService
}

func (d *diContainer) OrderPaidConsumer(ctx context.Context) service.OrderPaidConsumerService {
	if d.orderPaidConsumer == nil {
		notificationServiceInstance := d.NotificationService(ctx)

		d.orderPaidConsumer = consumer.NewOrderPaidService(
			d.OrderPaidKafkaConsumer(ctx),
			d.OrderPaidDecoder(ctx),
			notificationServiceInstance,
		)
	}
	return d.orderPaidConsumer
}

func (d *diContainer) ShipAssembledConsumer(ctx context.Context) service.ShipAssembledConsumerService {
	if d.shipAssembledConsumer == nil {
		notificationServiceInstance := d.NotificationService(ctx)

		d.shipAssembledConsumer = consumer.NewShipAssembledService(
			d.ShipAssembledKafkaConsumer(ctx),
			d.ShipAssembledDecoder(ctx),
			notificationServiceInstance,
		)
	}
	return d.shipAssembledConsumer
}

func (d *diContainer) IAMClient(ctx context.Context) iamClient.Client {
	if d.iamClient == nil {
		if os.Getenv("SKIP_GRPC_CONNECTIONS") == "true" {
			logger.Warn(ctx, "IAM connections skipped due to SKIP_GRPC_CONNECTIONS=true")
			return nil
		}

		grpcClient := d.IAMGRPCClient(ctx)
		if grpcClient != nil {
			d.iamClient = iamClient.NewClientWithGRPCClient(grpcClient)
		} else {
			client, err := iamClient.NewClient(ctx, config.AppConfig().IAM.IAMAddress())
			if err != nil {
				logger.Error(ctx, "Failed to create IAM client", zap.Error(err))
				panic(fmt.Sprintf("failed to create IAM client: %s\n", err.Error()))
			}

			closer.AddNamed("IAM client", func(ctx context.Context) error {
				return client.Close()
			})

			d.iamClient = client
		}
	}
	logger.Info(ctx, "IAM client created successfully")

	return d.iamClient
}

func (d *diContainer) IAMGRPCConn(ctx context.Context) *grpc.ClientConn {
	if d.iamGRPCConn == nil {
		if os.Getenv("SKIP_GRPC_CONNECTIONS") == "true" {
			return nil
		}

		conn, err := grpc.NewClient(
			config.AppConfig().IAM.IAMAddress(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			logger.Error(ctx, "Failed to create IAM gRPC connection", zap.Error(err))
			panic(fmt.Sprintf("failed to connect to IAM service: %v", err))
		}

		closer.AddNamed("IAM gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.iamGRPCConn = conn
	}
	return d.iamGRPCConn
}

func (d *diContainer) IAMGRPCClient(ctx context.Context) iamV1.UserServiceClient {
	if d.iamGRPCClient == nil {
		conn := d.IAMGRPCConn(ctx)
		if conn != nil {
			d.iamGRPCClient = iamV1.NewUserServiceClient(conn)
		}
	}
	return d.iamGRPCClient
}

func (d *diContainer) TelegramClient(ctx context.Context) telegramClient.TelegramClient {
	d.mu.RLock()
	if d.telegramClient != nil {
		defer d.mu.RUnlock()
		return d.telegramClient
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.telegramClient != nil {
		return d.telegramClient
	}

	d.telegramClientLocked = true
	client := d.telegramClientUnsafe(ctx)
	d.telegramClientLocked = false

	logger.Info(ctx, "Telegram client created successfully")

	return client
}

func (d *diContainer) telegramClientUnsafe(ctx context.Context) telegramClient.TelegramClient {
	if !d.telegramClientLocked {
		panic("telegramClientUnsafe должен вызываться под lock'ом")
	}

	if d.telegramClient != nil {
		return d.telegramClient
	}

	client, err := telegramClient.NewClient(ctx, config.AppConfig().Telegram)
	if err != nil {
		panic(fmt.Sprintf("failed to create telegram client: %s\n", err.Error()))
	}

	// Устанавливаем callback для регистрации пользователей сразу после создания клиента
	if d.iamClient != nil {
		// Создаем функцию-обертку для callback'а
		registrationCallback := func(ctx context.Context, username string, chatID int64) error {
			logger.Info(ctx, "Получен запрос на регистрацию пользователя из Telegram",
				zap.String("username", username),
				zap.Int64("chat_id", chatID))

			email := fmt.Sprintf("%s@telegram.local", username)
			password := fmt.Sprintf("Tg_%s_%d", username, chatID)

			notificationMethods := []*iamV1.NotificationMethod{
				{
					ProviderName: "telegram",
					Target:       strconv.FormatInt(chatID, 10),
				},
			}

			userUUID, err := d.iamClient.RegisterUser(ctx, username, email, password, notificationMethods)
			if err != nil {
				logger.Error(ctx, "Ошибка при вызове IAM сервиса для регистрации пользователя",
					zap.Error(err),
					zap.String("username", username),
					zap.String("email", email),
					zap.Int64("chat_id", chatID))
				return fmt.Errorf("failed to register user: %w", err)
			}

			logger.Info(ctx, "Пользователь успешно зарегистрирован через IAM сервис",
				zap.String("username", username),
				zap.String("email", email),
				zap.Int64("chat_id", chatID),
				zap.String("user_uuid", userUUID))

			return nil
		}

		client.SetUserRegistrationCallback(registrationCallback)
		logger.Info(ctx, "Установлен callback для регистрации пользователей в Telegram клиенте")
	} else {
		logger.Warn(ctx, "IAM клиент недоступен - callback для регистрации пользователей не установлен")
	}

	closer.AddNamed("Telegram client", func(ctx context.Context) error {
		return client.Close(ctx)
	})

	d.telegramClient = client
	return d.telegramClient
}

func (d *diContainer) OrderPaidConsumerGroup(ctx context.Context) sarama.ConsumerGroup {
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

	logger.Info(ctx, "OrderPaid Kafka consumer group created successfully")

	return d.orderPaidConsumerGroup
}

func (d *diContainer) ShipAssembledConsumerGroup(ctx context.Context) sarama.ConsumerGroup {
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

	logger.Info(ctx, "ShipAssembled Kafka consumer group created successfully")

	return d.shipAssembledConsumerGroup
}

func (d *diContainer) OrderPaidKafkaConsumer(ctx context.Context) wrappedKafka.Consumer {
	if d.orderPaidKafkaConsumer == nil {
		d.orderPaidKafkaConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderPaidConsumerGroup(ctx),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.Logger(),
		)
	}

	return d.orderPaidKafkaConsumer
}

func (d *diContainer) ShipAssembledKafkaConsumer(ctx context.Context) wrappedKafka.Consumer {
	if d.shipAssembledKafkaConsumer == nil {
		d.shipAssembledKafkaConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ShipAssembledConsumerGroup(ctx),
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
