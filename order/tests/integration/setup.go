//go:build integration

package integration

import (
	"context"
	"os"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/testcontainers"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/testcontainers/app"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/testcontainers/network"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/testcontainers/path"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/testcontainers/postgres"
)

const (
	// Параметры для контейнеров
	orderAppName    = "order-app"
	orderDockerfile = "deploy/docker/order/Dockerfile"

	// Переменные окружения приложения
	httpPortKey = "HTTP_PORT"

	// Значения переменных окружения
	loggerLevelValue = "debug"
	startupTimeout   = 3 * time.Minute
)

// TestEnvironment — структура для хранения ресурсов тестового окружения
type TestEnvironment struct {
	Network  *network.Network
	Postgres *postgres.Container
	App      *app.Container
}

// setupTestEnvironment — подготавливает тестовое окружение: сеть, контейнеры и возвращает структуру с ресурсами
func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Подготовка тестового окружения...")

	// Шаг 1: Создаём общую Docker-сеть
	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "не удалось создать общую сеть", zap.Error(err))
	}
	logger.Info(ctx, "✅ Сеть успешно создана")

	// Получаем переменные окружения для PostgreSQL с проверкой на наличие
	postgresUsername := getEnvWithLogging(ctx, testcontainers.PostgresUsernameKey)
	postgresPassword := getEnvWithLogging(ctx, testcontainers.PostgresPasswordKey)
	postgresImageName := getEnvWithLogging(ctx, testcontainers.PostgresImageNameKey)
	postgresDatabase := getEnvWithLogging(ctx, testcontainers.PostgresDatabaseKey)

	// Получаем порт HTTP для waitStrategy
	httpPort := getEnvWithLogging(ctx, httpPortKey)

	// Шаг 2: Запускаем контейнер с PostgreSQL
	generatedPostgres, err := postgres.NewContainer(ctx,
		postgres.WithNetworkName(generatedNetwork.Name()),
		postgres.WithContainerName(testcontainers.PostgresContainerName),
		postgres.WithImageName(postgresImageName),
		postgres.WithDatabase(postgresDatabase),
		postgres.WithAuth(postgresUsername, postgresPassword),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер PostgreSQL", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер PostgreSQL успешно запущен")

	// Шаг 3: Запускаем контейнер с приложением
	projectRoot := path.GetProjectRoot()

	// Создаем полный набор переменных окружения для приложения
	appEnv := map[string]string{
		// Настройки HTTP
		"HTTP_ADDRESS": "0.0.0.0:" + httpPort,

		// Настройки логгера
		"LOGGER_LEVEL":   "debug",
		"LOGGER_AS_JSON": "false",

		// Настройки PostgreSQL - используем значения из конфигурации PostgreSQL контейнера
		testcontainers.PostgresHostKey:     generatedPostgres.Config().ContainerName,
		testcontainers.PostgresPortKey:     "5432",
		testcontainers.PostgresDatabaseKey: generatedPostgres.Config().Database,
		testcontainers.PostgresUsernameKey: generatedPostgres.Config().Username,
		testcontainers.PostgresPasswordKey: generatedPostgres.Config().Password,
		"POSTGRES_SSLMODE":                 "disable",
		"POSTGRES_MIGRATIONS_DIR":          "migrations",
	}

	// Создаем настраиваемую стратегию ожидания с увеличенным таймаутом
	waitStrategy := wait.ForListeningPort(nat.Port(httpPort + "/tcp")).
		WithStartupTimeout(startupTimeout)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(orderAppName),
		app.WithPort(httpPort),
		app.WithDockerfile(projectRoot, orderDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Postgres: generatedPostgres})
		logger.Fatal(ctx, "не удалось запустить контейнер приложения", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер приложения успешно запущен")

	logger.Info(ctx, "🎉 Тестовое окружение готово")
	return &TestEnvironment{
		Network:  generatedNetwork,
		Postgres: generatedPostgres,
		App:      appContainer,
	}
}

// getEnvWithLogging возвращает значение переменной окружения с логированием
func getEnvWithLogging(ctx context.Context, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Warn(ctx, "Переменная окружения не установлена", zap.String("key", key))
	}

	return value
}
