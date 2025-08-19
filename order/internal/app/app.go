package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	customMiddleware "github.com/kont1n/MSA_Rocket_Factory/order/internal/api/middleware"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/config"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/closer"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	orderV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/openapi/order/v1"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
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
	// Запускаем Kafka Consumer в горутине только если не отключен для тестов
	if os.Getenv("SKIP_KAFKA_CONSUMER") != "true" {
		go func() {
			err := a.diContainer.ShipAssembledConsumer(ctx).RunConsumer(ctx)
			if err != nil {
				logger.Error(ctx, "❌ Ошибка при работе Kafka Consumer", zap.Error(err))
			}
		}()
	}

	return a.runHTTPServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initLogger,
		a.ensureDatabaseExists, // Проверяем и создаем БД перед инициализацией DI
		a.initDI,
		a.initCloser,
		a.initHTTPServer,
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

func (a *App) initHTTPServer(ctx context.Context) error {
	// Создаем OpenAPI сервер
	orderServer, err := orderV1.NewServer(a.diContainer.OrderV1API(ctx))
	if err != nil {
		return fmt.Errorf("failed to create OpenAPI server: %w", err)
	}

	// Настраиваем роутер
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(customMiddleware.RequestLogger)
	r.Mount("/", orderServer)

	// Создаем HTTP сервер
	a.httpServer = &http.Server{
		Addr:              config.AppConfig().HTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: time.Duration(config.AppConfig().HTTP.ReadHeaderTimeout()) * time.Second,
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(ctx, time.Duration(config.AppConfig().HTTP.ShutdownTimeout())*time.Second)
		defer cancel()

		err := a.httpServer.Shutdown(shutdownCtx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 HTTP Order service server listening on %s", config.AppConfig().HTTP.Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
