package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	grpcClient "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1API "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/api/v1"
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/config"
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository"
	inventoryRepository "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/mongo"
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/service"
	inventoryService "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/service/part"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/closer"
	grpcMiddleware "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/middleware/grpc"
	iamV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/iam/v1"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryAPIv1      inventoryV1.InventoryServiceServer
	inventoryService    service.InventoryService
	inventoryRepository repository.InventoryRepository
	mongoDBClient       *mongo.Client
	mongoDBHandle       *mongo.Database
	iamGRPCConn         *grpcClient.ClientConn
	iamClient           iamV1.AuthServiceClient
	authInterceptor     *grpcMiddleware.AuthInterceptor
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryAPIv1 == nil {
		d.inventoryAPIv1 = inventoryV1API.NewAPI(d.PartService(ctx))
	}
	return d.inventoryAPIv1
}

func (d *diContainer) PartService(ctx context.Context) service.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = inventoryService.NewService(d.PartRepository(ctx))
	}
	return d.inventoryService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.InventoryRepository {
	if d.inventoryRepository == nil {
		d.inventoryRepository = inventoryRepository.NewRepository(ctx, d.MongoDBHandle(ctx))
	}
	return d.inventoryRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}
	return d.mongoDBHandle
}

func (d *diContainer) IAMGRPCConn(ctx context.Context) *grpcClient.ClientConn {
	if d.iamGRPCConn == nil {
		// Для интеграционных тестов пропускаем подключение к реальному gRPC сервису
		if config.AppConfig().GRPCClient.IAMAddress() == "" {
			// Возвращаем nil - клиенты должны обрабатывать этот случай
			return nil
		}

		conn, err := grpcClient.NewClient(
			config.AppConfig().GRPCClient.IAMAddress(),
			grpcClient.WithTransportCredentials(insecure.NewCredentials()),
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

func (d *diContainer) IAMClient(ctx context.Context) iamV1.AuthServiceClient {
	if d.iamClient == nil {
		d.iamClient = iamV1.NewAuthServiceClient(d.IAMGRPCConn(ctx))
	}
	return d.iamClient
}

func (d *diContainer) AuthInterceptor(ctx context.Context) *grpcMiddleware.AuthInterceptor {
	if d.authInterceptor == nil {
		d.authInterceptor = grpcMiddleware.NewAuthInterceptor(d.IAMClient(ctx))
	}
	return d.authInterceptor
}
