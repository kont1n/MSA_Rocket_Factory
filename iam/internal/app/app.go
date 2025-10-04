package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/config"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/middleware"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/closer"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/grpc/health"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	iamV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/iam/v1"
	jwtV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/jwt/v1"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

// extractDBName извлекает имя базы данных из URI
func extractDBName(uri string) (string, error) {
	if len(uri) == 0 {
		return "", fmt.Errorf("пустой URI")
	}

	lastSlash := -1
	questionMark := -1
	for i := len(uri) - 1; i >= 0; i-- {
		if uri[i] == '?' && questionMark == -1 {
			questionMark = i
		}
		if uri[i] == '/' && lastSlash == -1 {
			lastSlash = i
			break
		}
	}

	if lastSlash == -1 || questionMark == -1 || lastSlash >= questionMark {
		return "", fmt.Errorf("неверный формат URI")
	}

	return uri[lastSlash+1 : questionMark], nil
}

// createSystemURI создает URI для подключения к системной БД postgres
func createSystemURI(uri string) (string, error) {
	if len(uri) == 0 {
		return "", fmt.Errorf("пустой URI")
	}

	lastSlash := -1
	questionMark := -1
	for i := len(uri) - 1; i >= 0; i-- {
		if uri[i] == '?' && questionMark == -1 {
			questionMark = i
		}
		if uri[i] == '/' && lastSlash == -1 {
			lastSlash = i
			break
		}
	}

	if lastSlash == -1 || questionMark == -1 || lastSlash >= questionMark {
		return "", fmt.Errorf("неверный формат URI")
	}

	return uri[:lastSlash+1] + "postgres" + uri[questionMark:], nil
}

// checkAndCreateDB проверяет существование БД и создает её при необходимости
func checkAndCreateDB(ctx context.Context, systemURI, dbName string) error {
	pool, err := pgxpool.New(ctx, systemURI)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к PostgreSQL: %w", err)
	}
	defer pool.Close()

	var exists bool
	err = pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка при проверке существования БД: %w", err)
	}

	if !exists {
		logger.Info(ctx, fmt.Sprintf("📝 База данных %s не существует, создаем...", dbName))
		_, err = pool.Exec(ctx, fmt.Sprintf("CREATE DATABASE \"%s\"", dbName))
		if err != nil {
			return fmt.Errorf("не удалось создать БД %s: %w", dbName, err)
		}
		logger.Info(ctx, fmt.Sprintf("✅ База данных %s успешно создана", dbName))
	} else {
		logger.Info(ctx, fmt.Sprintf("✅ База данных %s уже существует", dbName))
	}

	return nil
}

// ensureDatabaseExists проверяет существование базы данных и создает её при необходимости
func (a *App) ensureDatabaseExists(ctx context.Context) error {
	logger.Info(ctx, "🔍 Проверяем существование базы данных...")

	dbConfig := config.AppConfig().DB
	uri := dbConfig.URI()

	dbName, err := extractDBName(uri)
	if err != nil {
		return fmt.Errorf("не удалось извлечь имя БД из URI: %s", uri)
	}

	systemURI, err := createSystemURI(uri)
	if err != nil {
		return fmt.Errorf("не удалось создать URI для системной БД")
	}

	return checkAndCreateDB(ctx, systemURI, dbName)
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
		a.initLogger,
		a.ensureDatabaseExists, // Проверяем и создаем БД перед инициализацией DI
		a.initDI,
		a.initCloser,
		a.initMigration, // Выполняем миграции после инициализации DI
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
	grpcConfig := config.AppConfig().GRPC

	var opts []grpc.ServerOption

	// Добавляем rate limiting middleware (5 попыток входа в минуту)
	rateLimiter := middleware.NewRateLimiter(5, time.Minute)
	opts = append(opts, grpc.UnaryInterceptor(rateLimiter.UnaryServerInterceptor()))
	logger.Info(ctx, "✅ Rate limiting активирован: 5 попыток входа в минуту")

	if grpcConfig.IsInsecure() {
		// Режим разработки - использовать insecure соединения
		logger.Warn(ctx, "⚠️ gRPC сервер работает в небезопасном режиме. НЕ используйте в продакшене!")
		opts = append(opts, grpc.Creds(insecure.NewCredentials()))
	} else {
		// Продакшен - использовать TLS
		if grpcConfig.TLSCertFile() == "" || grpcConfig.TLSKeyFile() == "" {
			return fmt.Errorf("для продакшен-режима необходимы GRPC_TLS_CERT_FILE и GRPC_TLS_KEY_FILE")
		}

		creds, err := credentials.NewServerTLSFromFile(grpcConfig.TLSCertFile(), grpcConfig.TLSKeyFile())
		if err != nil {
			return fmt.Errorf("не удалось создать TLS credentials: %w", err)
		}

		opts = append(opts, grpc.Creds(creds))
		logger.Info(ctx, "✅ gRPC сервер использует TLS шифрование")
	}

	a.grpcServer = grpc.NewServer(opts...)
	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	// Регистрируем health service для проверки работоспособности
	health.RegisterService(a.grpcServer)

	iamV1.RegisterAuthServiceServer(a.grpcServer, a.diContainer.AuthV1API(ctx))
	iamV1.RegisterUserServiceServer(a.grpcServer, a.diContainer.UserV1API(ctx))
	jwtV1.RegisterJWTServiceServer(a.grpcServer, a.diContainer.JWTV1API(ctx))

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC IAM server listening on %s", config.AppConfig().GRPC.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}

// initMigration инициализирует и выполняет миграции базы данных
func (a *App) initMigration(ctx context.Context) error {
	logger.Info(ctx, "🔄 Начинаем выполнение миграций...")

	migrationsDir := config.AppConfig().DB.MigrationsDir()
	logger.Info(ctx, fmt.Sprintf("📁 Директория миграций: %s", migrationsDir))

	err := a.diContainer.MigrationRunner().Up()
	if err != nil {
		return fmt.Errorf("ошибка при выполнении миграций: %w", err)
	}

	logger.Info(ctx, "✅ Миграции успешно выполнены")
	return nil
}
