package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/kont1n/MSA_Rocket_Factory/order/internal/api/order/v1"
	grpcClients "github.com/kont1n/MSA_Rocket_Factory/order/internal/client/grpc"
	invClient "github.com/kont1n/MSA_Rocket_Factory/order/internal/client/grpc/inventory/v1"
	payClient "github.com/kont1n/MSA_Rocket_Factory/order/internal/client/grpc/payment/v1"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/config"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/repository"
	orderRepository "github.com/kont1n/MSA_Rocket_Factory/order/internal/repository/postgres"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/service"
	orderService "github.com/kont1n/MSA_Rocket_Factory/order/internal/service/order"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/closer"
	orderV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderAPIv1          orderV1.Handler
	orderService        service.OrderService
	orderRepository     repository.OrderRepository
	inventoryClient     grpcClients.InventoryClient
	paymentClient       grpcClients.PaymentClient
	dbPool              *pgxpool.Pool
	inventoryGRPCConn   *grpc.ClientConn
	paymentGRPCConn     *grpc.ClientConn
	inventoryGRPCClient inventoryV1.InventoryServiceClient
	paymentGRPCClient   paymentV1.PaymentServiceClient
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
		d.inventoryGRPCClient = inventoryV1.NewInventoryServiceClient(d.InventoryGRPCConn(ctx))
	}
	return d.inventoryGRPCClient
}

func (d *diContainer) PaymentGRPCClient(ctx context.Context) paymentV1.PaymentServiceClient {
	if d.paymentGRPCClient == nil {
		d.paymentGRPCClient = paymentV1.NewPaymentServiceClient(d.PaymentGRPCConn(ctx))
	}
	return d.paymentGRPCClient
}
