//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
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
	DBPool   *pgxpool.Pool
}

// setupTestEnvironment — подготавливает тестовое окружение: сеть, контейнеры и возвращает структуру с ресурсами
func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Подготовка тестового окружения...")

	// Генерируем уникальный суффикс для имен контейнеров, чтобы избежать конфликтов при параллельном запуске
	uniqueSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	uniqueNetworkName := fmt.Sprintf("%s-%s", projectName, uniqueSuffix)
	uniquePostgresName := fmt.Sprintf("%s-%s", testcontainers.PostgresContainerName, uniqueSuffix)

	// Шаг 1: Создаём общую Docker-сеть с уникальным именем
	generatedNetwork, err := network.NewNetwork(ctx, uniqueNetworkName)
	if err != nil {
		logger.Fatal(ctx, "не удалось создать общую сеть", zap.Error(err))
	}
	logger.Info(ctx, "✅ Сеть успешно создана", zap.String("network_name", uniqueNetworkName))

	// Получаем переменные окружения для PostgreSQL с проверкой на наличие и значениями по умолчанию
	postgresUsername := getEnvWithDefault(ctx, testcontainers.PostgresUsernameKey, testcontainers.PostgresUsername)
	postgresPassword := getEnvWithDefault(ctx, testcontainers.PostgresPasswordKey, testcontainers.PostgresPassword)
	postgresImageName := getEnvWithDefault(ctx, testcontainers.PostgresImageNameKey, testcontainers.PostgresImageName)
	postgresDatabase := getEnvWithDefault(ctx, testcontainers.PostgresDatabaseKey, testcontainers.PostgresDatabase)

	// Получаем порт HTTP для waitStrategy
	httpPort := getEnvWithDefault(ctx, httpPortKey, "8080")

	// Шаг 2: Запускаем контейнер с PostgreSQL с уникальным именем
	generatedPostgres, err := postgres.NewContainer(ctx,
		postgres.WithNetworkName(generatedNetwork.Name()),
		postgres.WithContainerName(uniquePostgresName),
		postgres.WithImageName(postgresImageName),
		postgres.WithDatabase(postgresDatabase),
		postgres.WithAuth(postgresUsername, postgresPassword),
		postgres.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер PostgreSQL", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер PostgreSQL успешно запущен")

	// Дополнительная задержка для полной готовности PostgreSQL
	logger.Info(ctx, "⏳ Ожидаем полной готовности PostgreSQL...")
	time.Sleep(5 * time.Second)

	// Дополнительная проверка готовности PostgreSQL
	logger.Info(ctx, "🔍 Проверяем готовность PostgreSQL...")
	testConnStr, err := generatedPostgres.ConnectionString(ctx)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Postgres: generatedPostgres})
		logger.Fatal(ctx, "не удалось получить строку подключения к PostgreSQL", zap.Error(err))
	}

	// Пытаемся подключиться и выполнить простой запрос
	testPool, err := pgxpool.New(ctx, testConnStr)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Postgres: generatedPostgres})
		logger.Fatal(ctx, "не удалось подключиться к PostgreSQL для проверки готовности", zap.Error(err))
	}
	defer testPool.Close()

	// Выполняем простой запрос для проверки готовности
	var result int
	err = testPool.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Postgres: generatedPostgres})
		logger.Fatal(ctx, "PostgreSQL не готов к выполнению запросов", zap.Error(err))
	}
	logger.Info(ctx, "✅ PostgreSQL полностью готов к работе")

	// Выполняем миграции базы данных
	logger.Info(ctx, "🔄 Выполняем миграции базы данных...")
	err = runMigrations(ctx, testConnStr)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Postgres: generatedPostgres})
		logger.Fatal(ctx, "не удалось выполнить миграции", zap.Error(err))
	}
	logger.Info(ctx, "✅ Миграции базы данных выполнены")

	// Шаг 3: Запускаем контейнер с приложением
	projectRoot := path.GetProjectRoot()

	// Создаем полный набор переменных окружения для приложения
	appEnv := map[string]string{
		// Настройки HTTP
		"HTTP_HOST": "0.0.0.0",
		"HTTP_PORT": httpPort,

		// Настройки логгера
		"LOGGER_LEVEL":   "debug",
		"LOGGER_AS_JSON": "false",

		// Настройки PostgreSQL - используем уникальное имя контейнера в сети
		"POSTGRES_HOST":           uniquePostgresName,
		"POSTGRES_PORT":           "5432",
		"POSTGRES_DATABASE":       generatedPostgres.Config().Database,
		"POSTGRES_USER":           generatedPostgres.Config().Username,
		"POSTGRES_PASSWORD":       generatedPostgres.Config().Password,
		"POSTGRES_SSLMODE":        "disable",
		"POSTGRES_MIGRATIONS_DIR": "migrations",

		// Настройки Kafka (фиктивные значения для интеграционных тестов)
		"KAFKA_BROKERS":       "localhost:9092",
		"CONSUMER_TOPIC_NAME": "ship.assembled",
		"CONSUMER_GROUP_ID":   "order-service",
		"PRODUCER_TOPIC_NAME": "order.paid",

		// Настройки gRPC клиентов (фиктивные значения для интеграционных тестов)
		"INVENTORY_GRPC_HOST": "localhost",
		"INVENTORY_GRPC_PORT": "50051",
		"PAYMENT_GRPC_HOST":   "localhost",
		"PAYMENT_GRPC_PORT":   "50052",

		// Отключаем Kafka Consumer для интеграционных тестов
		"SKIP_KAFKA_CONSUMER": "true",

		// Отключаем gRPC подключения для интеграционных тестов
		"SKIP_GRPC_CONNECTIONS": "true",

		// Отключаем проверку базы данных для интеграционных тестов
		"SKIP_DB_CHECK": "true",
	}

	// Создаем настраиваемую стратегию ожидания с увеличенным таймаутом
	// Ждем, что контейнер слушает порт
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

	// Дополнительная задержка для полной инициализации приложения
	logger.Info(ctx, "⏳ Ожидаем полной инициализации приложения...")
	time.Sleep(3 * time.Second)

	// Шаг 4: Создаем пул подключений к PostgreSQL
	connStr, err := generatedPostgres.ConnectionString(ctx)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Postgres: generatedPostgres, App: appContainer})
		logger.Fatal(ctx, "не удалось получить строку подключения к PostgreSQL", zap.Error(err))
	}

	dbPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Postgres: generatedPostgres, App: appContainer})
		logger.Fatal(ctx, "не удалось создать пул подключений к PostgreSQL", zap.Error(err))
	}
	logger.Info(ctx, "✅ Пул подключений к PostgreSQL создан")

	logger.Info(ctx, "🎉 Тестовое окружение готово")
	return &TestEnvironment{
		Network:  generatedNetwork,
		Postgres: generatedPostgres,
		App:      appContainer,
		DBPool:   dbPool,
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

// getEnvWithDefault возвращает значение переменной окружения или значение по умолчанию
func getEnvWithDefault(ctx context.Context, key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Info(ctx, "Используется значение по умолчанию для переменной окружения",
			zap.String("key", key),
			zap.String("default", defaultValue))
		return defaultValue
	}

	logger.Info(ctx, "Используется значение из переменной окружения",
		zap.String("key", key),
		zap.String("value", value))
	return value
}

// runMigrations выполняет миграции базы данных
func runMigrations(ctx context.Context, connStr string) error {
	// Подключаемся к базе данных
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к PostgreSQL для миграций: %w", err)
	}
	defer pool.Close()

	// Выполняем миграции с помощью goose
	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("не удалось установить диалект postgres: %w", err)
	}

	migrationsDir := path.GetProjectRoot() + "/order/migrations"
	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		return fmt.Errorf("не удалось выполнить миграции: %w", err)
	}

	return nil
}
