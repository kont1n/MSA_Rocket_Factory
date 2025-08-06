package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	inventoryV1API "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/api/v1"
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/config"
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository"
	inventoryRepository "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/mongo"
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/service"
	inventoryService "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/service/part"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/closer"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryAPIv1      inventoryV1.InventoryServiceServer
	inventoryService    service.InventoryService
	inventoryRepository repository.InventoryRepository
	mongoDBClient       *mongo.Client
	mongoDBHandle       *mongo.Database
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
		d.inventoryRepository = inventoryRepository.NewRepository(d.MongoDBHandle(ctx))
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
