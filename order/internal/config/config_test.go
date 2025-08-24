package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
	originalEnv map[string]string
}

func (s *ConfigSuite) SetupTest() {
	// Сохраняем текущие переменные окружения
	s.originalEnv = make(map[string]string)
	envVars := []string{
		"LOGGER_LEVEL",
		"LOGGER_AS_JSON",
		"HTTP_HOST",
		"HTTP_PORT",
		"HTTP_READ_HEADER_TIMEOUT",
		"HTTP_SHUTDOWN_TIMEOUT",
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_SSLMODE",
		"POSTGRES_DATABASE",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_MIGRATIONS_DIR",
		"INVENTORY_GRPC_HOST",
		"INVENTORY_GRPC_PORT",
		"PAYMENT_GRPC_HOST",
		"PAYMENT_GRPC_PORT",
		"KAFKA_BROKERS",
		"PRODUCER_TOPIC_NAME",
		"CONSUMER_TOPIC_NAME",
		"CONSUMER_GROUP_ID",
	}

	for _, envVar := range envVars {
		if val, exists := os.LookupEnv(envVar); exists {
			s.originalEnv[envVar] = val
		}
		_ = os.Unsetenv(envVar)
	}

	// Сбрасываем глобальную конфигурацию
	appConfig = nil
}

func (s *ConfigSuite) TearDownTest() {
	// Восстанавливаем оригинальные переменные окружения
	for envVar, val := range s.originalEnv {
		_ = os.Setenv(envVar, val)
	}

	// Очищаем переменные, которых не было
	envVars := []string{
		"LOGGER_LEVEL",
		"LOGGER_AS_JSON",
		"HTTP_HOST",
		"HTTP_PORT",
		"HTTP_READ_HEADER_TIMEOUT",
		"HTTP_SHUTDOWN_TIMEOUT",
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_SSLMODE",
		"POSTGRES_DATABASE",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_MIGRATIONS_DIR",
		"INVENTORY_GRPC_HOST",
		"INVENTORY_GRPC_PORT",
		"PAYMENT_GRPC_HOST",
		"PAYMENT_GRPC_PORT",
		"KAFKA_BROKERS",
		"PRODUCER_TOPIC_NAME",
		"CONSUMER_TOPIC_NAME",
		"CONSUMER_GROUP_ID",
	}

	for _, envVar := range envVars {
		if _, exists := s.originalEnv[envVar]; !exists {
			_ = os.Unsetenv(envVar)
		}
	}

	// Сбрасываем глобальную конфигурацию
	appConfig = nil
}

func (s *ConfigSuite) TestLoad_Success() {
	// Устанавливаем валидные переменные окружения
	_ = os.Setenv("LOGGER_LEVEL", "info")
	_ = os.Setenv("LOGGER_AS_JSON", "true")
	_ = os.Setenv("HTTP_HOST", "localhost")
	_ = os.Setenv("HTTP_PORT", "8080")
	_ = os.Setenv("HTTP_READ_HEADER_TIMEOUT", "10")
	_ = os.Setenv("HTTP_SHUTDOWN_TIMEOUT", "15")
	_ = os.Setenv("POSTGRES_HOST", "localhost")
	_ = os.Setenv("POSTGRES_PORT", "5432")
	_ = os.Setenv("POSTGRES_SSLMODE", "disable")
	_ = os.Setenv("POSTGRES_DATABASE", "orders")
	_ = os.Setenv("POSTGRES_USER", "user")
	_ = os.Setenv("POSTGRES_PASSWORD", "password")
	_ = os.Setenv("POSTGRES_MIGRATIONS_DIR", "./migrations")
	_ = os.Setenv("INVENTORY_GRPC_ADDRESS", "localhost:50051")
	_ = os.Setenv("PAYMENT_GRPC_ADDRESS", "localhost:50052")
	_ = os.Setenv("KAFKA_BROKERS", "localhost:9092")
	_ = os.Setenv("PRODUCER_TOPIC_NAME", "order-paid")
	_ = os.Setenv("CONSUMER_TOPIC_NAME", "ship-assembled")
	_ = os.Setenv("CONSUMER_GROUP_ID", "order-service")

	err := Load()
	s.NoError(err)
	s.NotNil(appConfig)

	// Проверяем, что конфигурация загружена корректно
	cfg := AppConfig()
	s.NotNil(cfg)
	s.Equal("info", cfg.Logger.Level())
	s.True(cfg.Logger.AsJson())
	s.Equal("localhost:8080", cfg.HTTP.Address())
	s.Equal(10, cfg.HTTP.ReadHeaderTimeout())
	s.Equal(15, cfg.HTTP.ShutdownTimeout())
	s.Equal("postgres://user:password@localhost:5432/orders?sslmode=disable", cfg.DB.URI())
	s.Equal("./migrations", cfg.DB.MigrationsDir())
	s.Equal("localhost:50051", cfg.GRPCClient.InventoryAddress())
	s.Equal("localhost:50052", cfg.GRPCClient.PaymentAddress())
}

func (s *ConfigSuite) TestLoad_HTTPConfigDefaults() {
	// Устанавливаем обязательные переменные, но не HTTP конфигурацию
	_ = os.Setenv("LOGGER_LEVEL", "info")
	_ = os.Setenv("LOGGER_AS_JSON", "true")
	_ = os.Setenv("POSTGRES_HOST", "localhost")
	_ = os.Setenv("POSTGRES_PORT", "5432")
	_ = os.Setenv("POSTGRES_SSLMODE", "disable")
	_ = os.Setenv("POSTGRES_DATABASE", "orders")
	_ = os.Setenv("POSTGRES_USER", "user")
	_ = os.Setenv("POSTGRES_PASSWORD", "password")
	_ = os.Setenv("POSTGRES_MIGRATIONS_DIR", "./migrations")
	_ = os.Setenv("KAFKA_BROKERS", "localhost:9092")
	_ = os.Setenv("PRODUCER_TOPIC_NAME", "order-paid")
	_ = os.Setenv("CONSUMER_TOPIC_NAME", "ship-assembled")
	_ = os.Setenv("CONSUMER_GROUP_ID", "order-service")

	err := Load()
	s.NoError(err)

	cfg := AppConfig()
	s.NotNil(cfg)
	// Проверяем значения по умолчанию для HTTP
	s.Equal("localhost:8080", cfg.HTTP.Address())
	s.Equal(5, cfg.HTTP.ReadHeaderTimeout())
	s.Equal(10, cfg.HTTP.ShutdownTimeout())
	// Проверяем значения по умолчанию для GRPC клиентов
	s.Equal("localhost:50051", cfg.GRPCClient.InventoryAddress())
	s.Equal("localhost:50052", cfg.GRPCClient.PaymentAddress())
}

func (s *ConfigSuite) TestLoad_MissingLoggerLevel() {
	// Устанавливаем все переменные кроме LOGGER_LEVEL
	_ = os.Setenv("LOGGER_AS_JSON", "true")
	_ = os.Setenv("POSTGRES_HOST", "localhost")
	_ = os.Setenv("POSTGRES_PORT", "5432")
	_ = os.Setenv("POSTGRES_SSLMODE", "disable")
	_ = os.Setenv("POSTGRES_DATABASE", "orders")
	_ = os.Setenv("POSTGRES_USER", "user")
	_ = os.Setenv("POSTGRES_PASSWORD", "password")
	_ = os.Setenv("POSTGRES_MIGRATIONS_DIR", "./migrations")
	_ = os.Setenv("KAFKA_BROKERS", "localhost:9092")
	_ = os.Setenv("PRODUCER_TOPIC_NAME", "order-paid")
	_ = os.Setenv("CONSUMER_TOPIC_NAME", "ship-assembled")
	_ = os.Setenv("CONSUMER_GROUP_ID", "order-service")

	err := Load()
	s.Error(err)
	s.Contains(err.Error(), "LOGGER_LEVEL")
}

func (s *ConfigSuite) TestLoad_MissingPostgresHost() {
	// Устанавливаем все переменные кроме POSTGRES_HOST
	_ = os.Setenv("LOGGER_LEVEL", "info")
	_ = os.Setenv("LOGGER_AS_JSON", "true")
	_ = os.Setenv("POSTGRES_PORT", "5432")
	_ = os.Setenv("POSTGRES_SSLMODE", "disable")
	_ = os.Setenv("POSTGRES_DATABASE", "orders")
	_ = os.Setenv("POSTGRES_USER", "user")
	_ = os.Setenv("POSTGRES_PASSWORD", "password")
	_ = os.Setenv("POSTGRES_MIGRATIONS_DIR", "./migrations")
	_ = os.Setenv("KAFKA_BROKERS", "localhost:9092")
	_ = os.Setenv("PRODUCER_TOPIC_NAME", "order-paid")
	_ = os.Setenv("CONSUMER_TOPIC_NAME", "ship-assembled")
	_ = os.Setenv("CONSUMER_GROUP_ID", "order-service")

	err := Load()
	s.Error(err)
	s.Contains(err.Error(), "POSTGRES_HOST")
}

func (s *ConfigSuite) TestLoad_MissingPostgresPassword() {
	// Устанавливаем все переменные кроме POSTGRES_PASSWORD
	_ = os.Setenv("LOGGER_LEVEL", "info")
	_ = os.Setenv("LOGGER_AS_JSON", "true")
	_ = os.Setenv("POSTGRES_HOST", "localhost")
	_ = os.Setenv("POSTGRES_PORT", "5432")
	_ = os.Setenv("POSTGRES_SSLMODE", "disable")
	_ = os.Setenv("POSTGRES_DATABASE", "orders")
	_ = os.Setenv("POSTGRES_USER", "user")
	_ = os.Setenv("POSTGRES_MIGRATIONS_DIR", "./migrations")
	_ = os.Setenv("KAFKA_BROKERS", "localhost:9092")
	_ = os.Setenv("PRODUCER_TOPIC_NAME", "order-paid")
	_ = os.Setenv("CONSUMER_TOPIC_NAME", "ship-assembled")
	_ = os.Setenv("CONSUMER_GROUP_ID", "order-service")

	err := Load()
	s.Error(err)
	s.Contains(err.Error(), "POSTGRES_PASSWORD")
}

func (s *ConfigSuite) TestLoad_FromEnvFile() {
	// Создаем временный .env файл
	tempDir := s.T().TempDir()
	envFile := filepath.Join(tempDir, ".env")

	envContent := `LOGGER_LEVEL=debug
LOGGER_AS_JSON=false
HTTP_HOST=0.0.0.0
HTTP_PORT=9090
HTTP_READ_HEADER_TIMEOUT=20
HTTP_SHUTDOWN_TIMEOUT=30
POSTGRES_HOST=db-host
POSTGRES_PORT=5433
POSTGRES_SSLMODE=require
POSTGRES_DATABASE=test_orders
POSTGRES_USER=testuser
POSTGRES_PASSWORD=testpass
POSTGRES_MIGRATIONS_DIR=/migrations
INVENTORY_GRPC_HOST=inventory
INVENTORY_GRPC_PORT=50051
PAYMENT_GRPC_HOST=payment
PAYMENT_GRPC_PORT=50052
KAFKA_BROKERS=localhost:9092
PRODUCER_TOPIC_NAME=order-paid
CONSUMER_TOPIC_NAME=ship-assembled
CONSUMER_GROUP_ID=order-service`

	err := os.WriteFile(envFile, []byte(envContent), 0o644)
	s.NoError(err)

	err = Load(envFile)
	s.NoError(err)

	cfg := AppConfig()
	s.NotNil(cfg)
	s.Equal("debug", cfg.Logger.Level())
	s.False(cfg.Logger.AsJson())
	s.Equal("0.0.0.0:9090", cfg.HTTP.Address())
	s.Equal(20, cfg.HTTP.ReadHeaderTimeout())
	s.Equal(30, cfg.HTTP.ShutdownTimeout())
	s.Equal("postgres://testuser:testpass@db-host:5433/test_orders?sslmode=require", cfg.DB.URI())
	s.Equal("/migrations", cfg.DB.MigrationsDir())
	s.Equal("inventory:50051", cfg.GRPCClient.InventoryAddress())
	s.Equal("payment:50052", cfg.GRPCClient.PaymentAddress())
}

func (s *ConfigSuite) TestLoad_NonExistentEnvFile() {
	// Тестируем загрузку с несуществующим файлом (не должно возвращать ошибку)
	// Но переменные окружения должны быть установлены
	_ = os.Setenv("LOGGER_LEVEL", "info")
	_ = os.Setenv("LOGGER_AS_JSON", "true")
	_ = os.Setenv("POSTGRES_HOST", "localhost")
	_ = os.Setenv("POSTGRES_PORT", "5432")
	_ = os.Setenv("POSTGRES_SSLMODE", "disable")
	_ = os.Setenv("POSTGRES_DATABASE", "orders")
	_ = os.Setenv("POSTGRES_USER", "user")
	_ = os.Setenv("POSTGRES_PASSWORD", "password")
	_ = os.Setenv("POSTGRES_MIGRATIONS_DIR", "./migrations")
	_ = os.Setenv("KAFKA_BROKERS", "localhost:9092")
	_ = os.Setenv("PRODUCER_TOPIC_NAME", "order-paid")
	_ = os.Setenv("CONSUMER_TOPIC_NAME", "ship-assembled")
	_ = os.Setenv("CONSUMER_GROUP_ID", "order-service")

	err := Load("/non/existent/file.env")
	s.NoError(err)

	cfg := AppConfig()
	s.NotNil(cfg)
}

func (s *ConfigSuite) TestAppConfig_BeforeLoad() {
	// Тестируем вызов AppConfig() до загрузки конфигурации
	cfg := AppConfig()
	s.Nil(cfg)
}

func (s *ConfigSuite) TestLoad_InvalidHTTPTimeout() {
	// Устанавливаем невалидное значение для HTTP_READ_HEADER_TIMEOUT
	_ = os.Setenv("LOGGER_LEVEL", "info")
	_ = os.Setenv("LOGGER_AS_JSON", "true")
	_ = os.Setenv("HTTP_READ_HEADER_TIMEOUT", "invalid")
	_ = os.Setenv("POSTGRES_HOST", "localhost")
	_ = os.Setenv("POSTGRES_PORT", "5432")
	_ = os.Setenv("POSTGRES_SSLMODE", "disable")
	_ = os.Setenv("POSTGRES_DATABASE", "orders")
	_ = os.Setenv("POSTGRES_USER", "user")
	_ = os.Setenv("POSTGRES_PASSWORD", "password")
	_ = os.Setenv("POSTGRES_MIGRATIONS_DIR", "./migrations")
	_ = os.Setenv("KAFKA_BROKERS", "localhost:9092")
	_ = os.Setenv("PRODUCER_TOPIC_NAME", "order-paid")
	_ = os.Setenv("CONSUMER_TOPIC_NAME", "ship-assembled")
	_ = os.Setenv("CONSUMER_GROUP_ID", "order-service")

	err := Load()
	s.NoError(err) // Невалидные значения заменяются на значения по умолчанию

	cfg := AppConfig()
	s.NotNil(cfg)
	s.Equal(5, cfg.HTTP.ReadHeaderTimeout()) // Значение по умолчанию
}

func (s *ConfigSuite) TestLoad_OverrideEnvFileWithEnvVars() {
	// Создаем .env файл с одними значениями
	tempDir := s.T().TempDir()
	envFile := filepath.Join(tempDir, ".env")

	envContent := `LOGGER_LEVEL=debug
LOGGER_AS_JSON=false
HTTP_HOST=localhost
HTTP_PORT=8080
POSTGRES_HOST=db-host
POSTGRES_PORT=5433
POSTGRES_SSLMODE=require
POSTGRES_DATABASE=test_orders
POSTGRES_USER=testuser
POSTGRES_PASSWORD=testpass
POSTGRES_MIGRATIONS_DIR=/migrations
KAFKA_BROKERS=localhost:9092
PRODUCER_TOPIC_NAME=order-paid
CONSUMER_TOPIC_NAME=ship-assembled
CONSUMER_GROUP_ID=order-service`

	err := os.WriteFile(envFile, []byte(envContent), 0o644)
	s.NoError(err)

	// Устанавливаем переменные окружения, которые должны переопределить значения из файла
	_ = os.Setenv("LOGGER_LEVEL", "error")
	_ = os.Setenv("POSTGRES_HOST", "override-host")
	_ = os.Setenv("HTTP_HOST", "override-host")
	_ = os.Setenv("HTTP_PORT", "3000")

	err = Load(envFile)
	s.NoError(err)

	cfg := AppConfig()
	s.NotNil(cfg)
	// Переменные окружения должны иметь приоритет над файлом
	s.Equal("error", cfg.Logger.Level())
	s.Equal("override-host:3000", cfg.HTTP.Address())
	s.Contains(cfg.DB.URI(), "override-host")
	// Остальные значения должны быть из файла
	s.False(cfg.Logger.AsJson())
	s.Equal("/migrations", cfg.DB.MigrationsDir())
}

func TestConfig(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
