package app

import (
	"context"
	"fmt"
	"os"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	customMiddleware "github.com/kont1n/MSA_Rocket_Factory/order/internal/api/middleware"
	orderV1API "github.com/kont1n/MSA_Rocket_Factory/order/internal/api/order/v1"
	grpcClients "github.com/kont1n/MSA_Rocket_Factory/order/internal/client/grpc"
	invClient "github.com/kont1n/MSA_Rocket_Factory/order/internal/client/grpc/inventory/v1"
	payClient "github.com/kont1n/MSA_Rocket_Factory/order/internal/client/grpc/payment/v1"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/config"
	kafkaConverter "github.com/kont1n/MSA_Rocket_Factory/order/internal/converter/kafka"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/converter/kafka/decoder"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/repository"
	orderRepository "github.com/kont1n/MSA_Rocket_Factory/order/internal/repository/postgres"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/service"
	shipAssembledConsumer "github.com/kont1n/MSA_Rocket_Factory/order/internal/service/consumer"
	orderService "github.com/kont1n/MSA_Rocket_Factory/order/internal/service/order"
	orderProducer "github.com/kont1n/MSA_Rocket_Factory/order/internal/service/producer"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/closer"
	wrappedKafka "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka/producer"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	orderV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/openapi/order/v1"
	iamV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/iam/v1"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderAPIv1                 orderV1.Handler
	orderService               service.OrderService
	orderRepository            repository.OrderRepository
	inventoryClient            grpcClients.InventoryClient
	paymentClient              grpcClients.PaymentClient
	dbPool                     *pgxpool.Pool
	inventoryGRPCConn          *grpc.ClientConn
	paymentGRPCConn            *grpc.ClientConn
	inventoryGRPCClient        inventoryV1.InventoryServiceClient
	paymentGRPCClient          paymentV1.PaymentServiceClient
	orderPaidProducer          service.OrderPaidProducer
	shipAssembledConsumer      service.ShipAssembledConsumer
	syncProducer               sarama.SyncProducer
	orderPaidKafkaProducer     wrappedKafka.Producer
	consumerGroup              sarama.ConsumerGroup
	shipAssembledKafkaConsumer wrappedKafka.Consumer
	shipAssembledDecoder       kafkaConverter.ShipAssembledDecoder
	iamGRPCConn                *grpc.ClientConn
	iamGRPCClient              iamV1.AuthServiceClient
	authMiddleware             *customMiddleware.AuthMiddleware
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderAPIv1 == nil {
		d.orderAPIv1 = orderV1API.NewAPI(d.OrderService(ctx))
	}
	return d.orderAPIv1
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
			d.OrderPaidProducer(ctx),
		)
	}
	return d.orderService
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(
			d.DBPool(ctx),
			config.AppConfig().DB.MigrationsDir(),
		)
	}
	return d.orderRepository
}

func (d *diContainer) InventoryClient(ctx context.Context) grpcClients.InventoryClient {
	if d.inventoryClient == nil {
		d.inventoryClient = invClient.NewClient(d.InventoryGRPCClient(ctx))
	}
	return d.inventoryClient
}

func (d *diContainer) PaymentClient(ctx context.Context) grpcClients.PaymentClient {
	if d.paymentClient == nil {
		d.paymentClient = payClient.NewClient(d.PaymentGRPCClient(ctx))
	}
	return d.paymentClient
}

func (d *diContainer) DBPool(ctx context.Context) *pgxpool.Pool {
	if d.dbPool == nil {
		// Для интеграционных тестов пропускаем подключение к реальной базе данных
		if os.Getenv("SKIP_DB_CHECK") == "true" {
			// Возвращаем nil - репозиторий должен обрабатывать этот случай
			return nil
		}

		pool, err := pgxpool.New(ctx, config.AppConfig().DB.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to connect to database: %v", err))
		}

		closer.AddNamed("DB pool", func(ctx context.Context) error {
			d.dbPool.Close()
			return nil
		})

		d.dbPool = pool
	}
	return d.dbPool
}

func (d *diContainer) InventoryGRPCConn(ctx context.Context) *grpc.ClientConn {
	if d.inventoryGRPCConn == nil {
		// Для интеграционных тестов пропускаем подключение к реальному gRPC сервису
		if os.Getenv("SKIP_GRPC_CONNECTIONS") == "true" {
			// Возвращаем nil - клиенты должны обрабатывать этот случай
			return nil
		}

		conn, err := grpc.NewClient(
			config.AppConfig().GRPCClient.InventoryAddress(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to inventory service: %v", err))
		}

		closer.AddNamed("Inventory gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.inventoryGRPCConn = conn
	}
	return d.inventoryGRPCConn
}

func (d *diContainer) PaymentGRPCConn(ctx context.Context) *grpc.ClientConn {
	if d.paymentGRPCConn == nil {
		// Для интеграционных тестов пропускаем подключение к реальному gRPC сервису
		if os.Getenv("SKIP_GRPC_CONNECTIONS") == "true" {
			// Возвращаем nil - клиенты должны обрабатывать этот случай
			return nil
		}

		conn, err := grpc.NewClient(
			config.AppConfig().GRPCClient.PaymentAddress(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to payment service: %v", err))
		}

		closer.AddNamed("Payment gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.paymentGRPCConn = conn
	}
	return d.paymentGRPCConn
}

func (d *diContainer) InventoryGRPCClient(ctx context.Context) inventoryV1.InventoryServiceClient {
	if d.inventoryGRPCClient == nil {
		conn := d.InventoryGRPCConn(ctx)
		if conn != nil {
			d.inventoryGRPCClient = inventoryV1.NewInventoryServiceClient(conn)
		}
		// Для интеграционных тестов возвращаем nil - клиенты должны обрабатывать этот случай
	}
	return d.inventoryGRPCClient
}

func (d *diContainer) PaymentGRPCClient(ctx context.Context) paymentV1.PaymentServiceClient {
	if d.paymentGRPCClient == nil {
		conn := d.PaymentGRPCConn(ctx)
		if conn != nil {
			d.paymentGRPCClient = paymentV1.NewPaymentServiceClient(conn)
		}
		// Для интеграционных тестов возвращаем nil - клиенты должны обрабатывать этот случай
	}
	return d.paymentGRPCClient
}

func (d *diContainer) OrderPaidProducer(ctx context.Context) service.OrderPaidProducer {
	if d.orderPaidProducer == nil {
		if os.Getenv("SKIP_KAFKA_CONSUMER") == "true" {
			// Для интеграционных тестов возвращаем nil
			return nil
		}
		d.orderPaidProducer = orderProducer.NewService(d.OrderPaidKafkaProducer())
	}
	return d.orderPaidProducer
}

func (d *diContainer) ShipAssembledConsumer(ctx context.Context) service.ShipAssembledConsumer {
	if d.shipAssembledConsumer == nil {
		if os.Getenv("SKIP_KAFKA_CONSUMER") == "true" {
			// Для интеграционных тестов возвращаем nil
			return nil
		}
		d.shipAssembledConsumer = shipAssembledConsumer.NewService(
			d.ShipAssembledKafkaConsumer(),
			d.ShipAssembledDecoder(ctx),
			d.OrderService(ctx),
		)
	}
	return d.shipAssembledConsumer
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidProducer.Config(),
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

func (d *diContainer) OrderPaidKafkaProducer() wrappedKafka.Producer {
	if d.orderPaidKafkaProducer == nil {
		d.orderPaidKafkaProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.orderPaidKafkaProducer
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().ShipAssembledConsumer.GroupID(),
			config.AppConfig().ShipAssembledConsumer.Config(),
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

func (d *diContainer) ShipAssembledKafkaConsumer() wrappedKafka.Consumer {
	if d.shipAssembledKafkaConsumer == nil {
		d.shipAssembledKafkaConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().ShipAssembledConsumer.Topic(),
			},
			logger.Logger(),
		)
	}

	return d.shipAssembledKafkaConsumer
}

func (d *diContainer) ShipAssembledDecoder(ctx context.Context) kafkaConverter.ShipAssembledDecoder {
	if d.shipAssembledDecoder == nil {
		d.shipAssembledDecoder = decoder.NewShipAssembledDecoder()
	}

	return d.shipAssembledDecoder
}

func (d *diContainer) IAMGRPCConn(ctx context.Context) *grpc.ClientConn {
	if d.iamGRPCConn == nil {
		// Для интеграционных тестов пропускаем подключение к реальному gRPC сервису
		if os.Getenv("SKIP_GRPC_CONNECTIONS") == "true" {
			// Возвращаем nil - клиенты должны обрабатывать этот случай
			return nil
		}

		conn, err := grpc.NewClient(
			config.AppConfig().GRPCClient.IAMAddress(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to IAM service: %v", err))
		}

		closer.AddNamed("IAM gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.iamGRPCConn = conn
	}
	return d.iamGRPCConn
}

func (d *diContainer) IAMGRPCClient(ctx context.Context) iamV1.AuthServiceClient {
	if d.iamGRPCClient == nil {
		conn := d.IAMGRPCConn(ctx)
		if conn != nil {
			d.iamGRPCClient = iamV1.NewAuthServiceClient(conn)
		}
		// Для интеграционных тестов возвращаем nil - клиенты должны обрабатывать этот случай
	}
	return d.iamGRPCClient
}

func (d *diContainer) AuthMiddleware(ctx context.Context) *customMiddleware.AuthMiddleware {
	if d.authMiddleware == nil {
		d.authMiddleware = customMiddleware.NewAuthMiddleware(d.IAMGRPCClient(ctx))
	}
	return d.authMiddleware
}
