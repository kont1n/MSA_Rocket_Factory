package postgres

import (
	"context"

	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/testcontainers"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

// Config представляет конфигурацию PostgreSQL контейнера
type Config struct {
	NetworkName   string
	ContainerName string
	ImageName     string
	Database      string
	Username      string
	Password      string
	Port          string
	Logger        Logger
}

// NewConfig создает новую конфигурацию PostgreSQL с значениями по умолчанию
func NewConfig() *Config {
	return &Config{
		NetworkName:   "test-network",
		ContainerName: testcontainers.PostgresContainerName,
		ImageName:     testcontainers.PostgresImageName,
		Database:      testcontainers.PostgresDatabase,
		Username:      testcontainers.PostgresUsername,
		Password:      testcontainers.PostgresPassword,
		Port:          testcontainers.PostgresPort,
		Logger:        &logger.NoopLogger{},
	}
}
