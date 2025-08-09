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
		"HTTP_HOST",
		"HTTP_PORT",
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
		"HTTP_HOST",
		"HTTP_PORT",
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
	_ = os.Setenv("GRPC_PORT", "50052")
	_ = os.Setenv("HTTP_HOST", "0.0.0.0")
	_ = os.Setenv("HTTP_PORT", "8080")

	err := Load()
	s.NoError(err)
	s.NotNil(appConfig)

	// Проверяем, что конфигурация загружена корректно
	cfg := AppConfig()
	s.NotNil(cfg)
	s.Equal("info", cfg.Logger.Level())
	s.True(cfg.Logger.AsJson())
	s.Equal("localhost:50052", cfg.GRPC.Address())
	s.Equal("0.0.0.0:8080", cfg.Http.Address())
}

func (s *ConfigSuite) TestLoad_MissingLoggerLevel() {
	// Устанавливаем все переменные кроме LOGGER_LEVEL
	_ = os.Setenv("LOGGER_AS_JSON", "true")
	_ = os.Setenv("GRPC_HOST", "localhost")
	_ = os.Setenv("GRPC_PORT", "50052")
	_ = os.Setenv("HTTP_HOST", "0.0.0.0")
	_ = os.Setenv("HTTP_PORT", "8080")

	err := Load()
	s.Error(err)
	s.Contains(err.Error(), "LOGGER_LEVEL")
}

func (s *ConfigSuite) TestLoad_MissingGRPCHost() {
	// Устанавливаем все переменные кроме GRPC_HOST
	_ = os.Setenv("LOGGER_LEVEL", "info")
	_ = os.Setenv("LOGGER_AS_JSON", "true")
	_ = os.Setenv("GRPC_PORT", "50052")
	_ = os.Setenv("HTTP_HOST", "0.0.0.0")
	_ = os.Setenv("HTTP_PORT", "8080")

	err := Load()
	s.Error(err)
	s.Contains(err.Error(), "GRPC_HOST")
}

func (s *ConfigSuite) TestLoad_MissingHTTPHost() {
	// Устанавливаем все переменные кроме HTTP_HOST
	_ = os.Setenv("LOGGER_LEVEL", "info")
	_ = os.Setenv("LOGGER_AS_JSON", "true")
	_ = os.Setenv("GRPC_HOST", "localhost")
	_ = os.Setenv("GRPC_PORT", "50052")
	_ = os.Setenv("HTTP_PORT", "8080")

	err := Load()
	s.Error(err)
	s.Contains(err.Error(), "HTTP_HOST")
}

func (s *ConfigSuite) TestLoad_MissingHTTPPort() {
	// Устанавливаем все переменные кроме HTTP_PORT
	_ = os.Setenv("LOGGER_LEVEL", "info")
	_ = os.Setenv("LOGGER_AS_JSON", "true")
	_ = os.Setenv("GRPC_HOST", "localhost")
	_ = os.Setenv("GRPC_PORT", "50052")
	_ = os.Setenv("HTTP_HOST", "0.0.0.0")

	err := Load()
	s.Error(err)
	s.Contains(err.Error(), "HTTP_PORT")
}

func (s *ConfigSuite) TestLoad_FromEnvFile() {
	// Создаем временный .env файл
	tempDir := s.T().TempDir()
	envFile := filepath.Join(tempDir, ".env")

	envContent := `LOGGER_LEVEL=debug
LOGGER_AS_JSON=false
GRPC_HOST=0.0.0.0
GRPC_PORT=9090
HTTP_HOST=127.0.0.1
HTTP_PORT=3000`

	err := os.WriteFile(envFile, []byte(envContent), 0o644)
	s.NoError(err)

	err = Load(envFile)
	s.NoError(err)

	cfg := AppConfig()
	s.NotNil(cfg)
	s.Equal("debug", cfg.Logger.Level())
	s.False(cfg.Logger.AsJson())
	s.Equal("0.0.0.0:9090", cfg.GRPC.Address())
	s.Equal("127.0.0.1:3000", cfg.Http.Address())
}

func (s *ConfigSuite) TestLoad_NonExistentEnvFile() {
	// Тестируем загрузку с несуществующим файлом (не должно возвращать ошибку)
	// Но переменные окружения должны быть установлены
	_ = os.Setenv("LOGGER_LEVEL", "info")
	_ = os.Setenv("LOGGER_AS_JSON", "true")
	_ = os.Setenv("GRPC_HOST", "localhost")
	_ = os.Setenv("GRPC_PORT", "50052")
	_ = os.Setenv("HTTP_HOST", "0.0.0.0")
	_ = os.Setenv("HTTP_PORT", "8080")

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
	_ = os.Setenv("GRPC_PORT", "50052")
	_ = os.Setenv("HTTP_HOST", "0.0.0.0")
	_ = os.Setenv("HTTP_PORT", "8080")

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
HTTP_HOST=127.0.0.1
HTTP_PORT=3000`

	err := os.WriteFile(envFile, []byte(envContent), 0o644)
	s.NoError(err)

	// Устанавливаем переменные окружения, которые должны переопределить значения из файла
	_ = os.Setenv("LOGGER_LEVEL", "error")
	_ = os.Setenv("GRPC_PORT", "8080")
	_ = os.Setenv("HTTP_PORT", "9000")

	err = Load(envFile)
	s.NoError(err)

	cfg := AppConfig()
	s.NotNil(cfg)
	// Переменные окружения должны иметь приоритет над файлом
	s.Equal("error", cfg.Logger.Level())
	s.Equal("0.0.0.0:8080", cfg.GRPC.Address())
	s.Equal("127.0.0.1:9000", cfg.Http.Address())
	// Остальные значения должны быть из файла
	s.False(cfg.Logger.AsJson())
}

func (s *ConfigSuite) TestLoad_DifferentLoggerLevels() {
	// Тестируем различные уровни логирования
	testCases := []string{"debug", "info", "warn", "error"}

	for _, level := range testCases {
		s.Run("LoggerLevel_"+level, func() {
			// Очищаем конфигурацию перед каждым тестом
			appConfig = nil

			_ = os.Setenv("LOGGER_LEVEL", level)
			_ = os.Setenv("LOGGER_AS_JSON", "false")
			_ = os.Setenv("GRPC_HOST", "localhost")
			_ = os.Setenv("GRPC_PORT", "50052")
			_ = os.Setenv("HTTP_HOST", "0.0.0.0")
			_ = os.Setenv("HTTP_PORT", "8080")

			err := Load()
			s.NoError(err)

			cfg := AppConfig()
			s.NotNil(cfg)
			s.Equal(level, cfg.Logger.Level())
		})
	}
}

func (s *ConfigSuite) TestLoad_BooleanLoggerAsJson() {
	// Тестируем различные булевые значения для LOGGER_AS_JSON
	testCases := []struct {
		value    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1", true},
		{"0", false},
		{"True", true},
		{"False", false},
		{"TRUE", true},
		{"FALSE", false},
	}

	for _, tc := range testCases {
		s.Run("LoggerAsJson_"+tc.value, func() {
			// Очищаем конфигурацию перед каждым тестом
			appConfig = nil

			_ = os.Setenv("LOGGER_LEVEL", "info")
			_ = os.Setenv("LOGGER_AS_JSON", tc.value)
			_ = os.Setenv("GRPC_HOST", "localhost")
			_ = os.Setenv("GRPC_PORT", "50052")
			_ = os.Setenv("HTTP_HOST", "0.0.0.0")
			_ = os.Setenv("HTTP_PORT", "8080")

			err := Load()
			s.NoError(err)

			cfg := AppConfig()
			s.NotNil(cfg)
			s.Equal(tc.expected, cfg.Logger.AsJson())
		})
	}
}

func TestConfigIntegration(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
