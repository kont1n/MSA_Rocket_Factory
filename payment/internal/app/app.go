package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/config"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/closer"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/grpc/health"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	paymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
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
	return a.runServers(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initGRPCServer,
		a.initGateway,
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

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
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
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	// Регистрируем health service для проверки работоспособности
	health.RegisterService(a.grpcServer)

	paymentV1.RegisterPaymentServiceServer(a.grpcServer, a.diContainer.PaymentV1API(ctx))

	return nil
}

func (a *App) initGateway(ctx context.Context) error {
	gateway := a.diContainer.Gateway(ctx)

	// Добавляем gateway в closer для корректного завершения
	closer.AddNamed("HTTP Gateway", func(ctx context.Context) error {
		return gateway.Stop(ctx)
	})

	return gateway.RegisterHandlers(ctx)
}

func (a *App) runServers(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	// Запускаем gRPC сервер
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info(ctx, fmt.Sprintf("🚀 gRPC PaymentService server listening on %s", config.AppConfig().GRPC.Address()))

		err := a.grpcServer.Serve(a.listener)
		if err != nil {
			errCh <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	// Запускаем HTTP Gateway
	wg.Add(1)
	go func() {
		defer wg.Done()
		gateway := a.diContainer.Gateway(ctx)

		err := gateway.Start(ctx)
		if err != nil {
			errCh <- fmt.Errorf("HTTP gateway error: %w", err)
		}
	}()

	// Ждем завершения или ошибки
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Возвращаем первую ошибку, если она есть
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}
