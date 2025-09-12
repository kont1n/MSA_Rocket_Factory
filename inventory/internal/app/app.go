package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/api/health"
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/config"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/closer"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/tracing"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initTracing,
		a.initCloser,
		a.initListener,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	return logger.Init(
		ctx,
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
		config.AppConfig().Logger.Outputs(),
		config.AppConfig().Logger.OtelEndpoint(),
		config.AppConfig().Logger.ServiceName(),
		"dev", // serviceEnvironment - хардкод для обратной совместимости
	)
}

func (a *App) initTracing(ctx context.Context) error {
	return tracing.InitTracer(ctx, config.AppConfig().Tracing)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())

	// Добавляем graceful shutdown для tracing
	closer.AddNamed("Tracing provider", tracing.ShutdownTracer)

	return nil
}

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().GRPC.Address())
	if err != nil {
		return err
	}
	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})

	a.listener = listener

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	// Создаем gRPC сервер с интерцепторами
	authInterceptor := a.diContainer.AuthInterceptor(ctx)

	var opts []grpc.ServerOption
	opts = append(opts, grpc.Creds(insecure.NewCredentials()))

	// Создаем цепочку интерцепторов
	var interceptors []grpc.UnaryServerInterceptor

	// Добавляем трейсинг интерцептор первым
	interceptors = append(interceptors, tracing.UnaryServerInterceptor("inventory"))

	// Добавляем AuthInterceptor если он создан
	if authInterceptor != nil {
		interceptors = append(interceptors, authInterceptor.Unary())
	}

	// Применяем цепочку интерцепторов
	if len(interceptors) > 0 {
		opts = append(opts, grpc.ChainUnaryInterceptor(interceptors...))
	}

	a.grpcServer = grpc.NewServer(opts...)

	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	// Регистрируем health service для проверки работоспособности (без аутентификации)
	// Health endpoint доступен через reflection API
	health.RegisterService(a.grpcServer)

	// Регистрируем основные API с аутентификацией
	inventoryV1.RegisterInventoryServiceServer(a.grpcServer, a.diContainer.InventoryV1API(ctx))

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC InventoryService server listening on %s", config.AppConfig().GRPC.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}
