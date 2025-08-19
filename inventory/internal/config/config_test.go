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
		"GRPC_HOST",
		"GRPC_PORT",
		"MONGO_HOST",
		"MONGO_PORT",
		"MONGO_DATABASE",
		"MONGO_INITDB_ROOT_USERNAME",
		"MONGO_INITDB_ROOT_PASSWORD",
		"MONGO_AUTH_DB",
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
		"GRPC_HOST",
		"GRPC_PORT",
		"MONGO_HOST",
		"MONGO_PORT",
		"MONGO_DATABASE",
		"MONGO_INITDB_ROOT_USERNAME",
		"MONGO_INITDB_ROOT_PASSWORD",
		"MONGO_AUTH_DB",
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
	_ = os.Setenv("GRPC_HOST", "localhost")
	_ = os.Setenv("GRPC_PORT", "50051")
	_ = os.Setenv("MONGO_HOST", "localhost")
	_ = os.Setenv("MONGO_PORT", "27017")
	_ = os.Setenv("MONGO_DATABASE", "inventory")
	_ = os.Setenv("MONGO_INITDB_ROOT_USERNAME", "admin")
	_ = os.Setenv("MONGO_INITDB_ROOT_PASSWORD", "password")
	_ = os.Setenv("MONGO_AUTH_DB", "admin")

	err := Load()
	s.NoError(err)
	s.NotNil(appConfig)

	// Проверяем, что конфигурация загружена корректно
	cfg := AppConfig()
	s.NotNil(cfg)
	s.Equal("info", cfg.Logger.Level())
	s.True(cfg.Logger.AsJson())
	s.Equal("localhost:50051", cfg.GRPC.Address())
	s.Equal("mongodb://admin:password@localhost:27017/inventory?authSource=admin", cfg.Mongo.URI())
	s.Equal("inventory", cfg.Mongo.DatabaseName())
}

func (s *ConfigSuite) TestLoad_MissingLoggerLevel() {
	// Устанавливаем все переменные кроме LOGGER_LEVEL
	_ = os.Setenv("LOGGER_AS_JSON", "true")
	_ = os.Setenv("GRPC_HOST", "localhost")
	_ = os.Setenv("GRPC_PORT", "50051")
	_ = os.Setenv("MONGO_HOST", "localhost")
	_ = os.Setenv("MONGO_PORT", "27017")
	_ = os.Setenv("MONGO_DATABASE", "inventory")
	_ = os.Setenv("MONGO_INITDB_ROOT_USERNAME", "admin")
	_ = os.Setenv("MONGO_INITDB_ROOT_PASSWORD", "password")
	_ = os.Setenv("MONGO_AUTH_DB", "admin")

	err := Load()
	s.Error(err)
	s.Contains(err.Error(), "LOGGER_LEVEL")
}

func (s *ConfigSuite) TestLoad_MissingGRPCHost() {
	// Устанавливаем все переменные кроме GRPC_HOST
	_ = os.Setenv("LOGGER_LEVEL", "info")
	_ = os.Setenv("LOGGER_AS_JSON", "true")
	_ = os.Setenv("GRPC_PORT", "50051")
	_ = os.Setenv("MONGO_HOST", "localhost")
	_ = os.Setenv("MONGO_PORT", "27017")
	_ = os.Setenv("MONGO_DATABASE", "inventory")
	_ = os.Setenv("MONGO_INITDB_ROOT_USERNAME", "admin")
	_ = os.Setenv("MONGO_INITDB_ROOT_PASSWORD", "password")
	_ = os.Setenv("MONGO_AUTH_DB", "admin")

	err := Load()
	s.Error(err)
	s.Contains(err.Error(), "GRPC_HOST")
}

func (s *ConfigSuite) TestLoad_MissingMongoHost() {
	// Устанавливаем все переменные кроме MONGO_HOST
	_ = os.Setenv("LOGGER_LEVEL", "info")
	_ = os.Setenv("LOGGER_AS_JSON", "true")
	_ = os.Setenv("GRPC_HOST", "localhost")
	_ = os.Setenv("GRPC_PORT", "50051")
	_ = os.Setenv("MONGO_PORT", "27017")
	_ = os.Setenv("MONGO_DATABASE", "inventory")
	_ = os.Setenv("MONGO_INITDB_ROOT_USERNAME", "admin")
	_ = os.Setenv("MONGO_INITDB_ROOT_PASSWORD", "password")
	_ = os.Setenv("MONGO_AUTH_DB", "admin")

	err := Load()
	s.Error(err)
	s.Contains(err.Error(), "MONGO_HOST")
}

func (s *ConfigSuite) TestLoad_FromEnvFile() {
	// Создаем временный .env файл
	tempDir := s.T().TempDir()
	envFile := filepath.Join(tempDir, ".env")

	envContent := `LOGGER_LEVEL=debug
LOGGER_AS_JSON=false
GRPC_HOST=0.0.0.0
GRPC_PORT=9090
MONGO_HOST=test-host
MONGO_PORT=27018
MONGO_DATABASE=test_inventory
MONGO_INITDB_ROOT_USERNAME=testuser
MONGO_INITDB_ROOT_PASSWORD=testpass
MONGO_AUTH_DB=admin`

	err := os.WriteFile(envFile, []byte(envContent), 0o644)
	s.NoError(err)

	err = Load(envFile)
	s.NoError(err)

	cfg := AppConfig()
	s.NotNil(cfg)
	s.Equal("debug", cfg.Logger.Level())
	s.False(cfg.Logger.AsJson())
	s.Equal("0.0.0.0:9090", cfg.GRPC.Address())
	s.Equal("mongodb://testuser:testpass@test-host:27018/test_inventory?authSource=admin", cfg.Mongo.URI())
	s.Equal("test_inventory", cfg.Mongo.DatabaseName())
}

func (s *ConfigSuite) TestLoad_NonExistentEnvFile() {
	// Тестируем загрузку с несуществующим файлом (не должно возвращать ошибку)
	// Но переменные окружения должны быть установлены
	_ = os.Setenv("LOGGER_LEVEL", "info")
	_ = os.Setenv("LOGGER_AS_JSON", "true")
	_ = os.Setenv("GRPC_HOST", "localhost")
	_ = os.Setenv("GRPC_PORT", "50051")
	_ = os.Setenv("MONGO_HOST", "localhost")
	_ = os.Setenv("MONGO_PORT", "27017")
	_ = os.Setenv("MONGO_DATABASE", "inventory")
	_ = os.Setenv("MONGO_INITDB_ROOT_USERNAME", "admin")
	_ = os.Setenv("MONGO_INITDB_ROOT_PASSWORD", "password")
	_ = os.Setenv("MONGO_AUTH_DB", "admin")

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

func (s *ConfigSuite) TestLoad_InvalidLoggerAsJson() {
	// Устанавливаем невалидное значение для LOGGER_AS_JSON
	_ = os.Setenv("LOGGER_LEVEL", "info")
	_ = os.Setenv("LOGGER_AS_JSON", "invalid")
	_ = os.Setenv("GRPC_HOST", "localhost")
	_ = os.Setenv("GRPC_PORT", "50051")
	_ = os.Setenv("MONGO_HOST", "localhost")
	_ = os.Setenv("MONGO_PORT", "27017")
	_ = os.Setenv("MONGO_DATABASE", "inventory")
	_ = os.Setenv("MONGO_INITDB_ROOT_USERNAME", "admin")
	_ = os.Setenv("MONGO_INITDB_ROOT_PASSWORD", "password")
	_ = os.Setenv("MONGO_AUTH_DB", "admin")

	err := Load()
	s.Error(err)
}

func (s *ConfigSuite) TestLoad_OverrideEnvFileWithEnvVars() {
	// Создаем .env файл с одними значениями
	tempDir := s.T().TempDir()
	envFile := filepath.Join(tempDir, ".env")

	envContent := `LOGGER_LEVEL=debug
LOGGER_AS_JSON=false
GRPC_HOST=0.0.0.0
GRPC_PORT=9090
MONGO_HOST=test-host
MONGO_PORT=27018
MONGO_DATABASE=test_inventory
MONGO_INITDB_ROOT_USERNAME=testuser
MONGO_INITDB_ROOT_PASSWORD=testpass
MONGO_AUTH_DB=admin`

	err := os.WriteFile(envFile, []byte(envContent), 0o644)
	s.NoError(err)

	// Устанавливаем переменные окружения, которые должны переопределить значения из файла
	_ = os.Setenv("LOGGER_LEVEL", "error")
	_ = os.Setenv("GRPC_PORT", "8080")
	_ = os.Setenv("MONGO_HOST", "override-host")
	_ = os.Setenv("MONGO_PORT", "8080")

	err = Load(envFile)
	s.NoError(err)

	cfg := AppConfig()
	s.NotNil(cfg)
	// Переменные окружения должны иметь приоритет над файлом
	s.Equal("error", cfg.Logger.Level())
	s.Equal("0.0.0.0:8080", cfg.GRPC.Address())
	// Остальные значения должны быть из файла
	s.False(cfg.Logger.AsJson())
	s.Equal("mongodb://testuser:testpass@override-host:8080/test_inventory?authSource=admin", cfg.Mongo.URI())
}

func TestConfig(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
