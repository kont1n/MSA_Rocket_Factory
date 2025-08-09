package mongo

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Config struct {
	NetworkName   string
	ContainerName string
	ImageName     string
	Database      string
	Username      string
	Password      string
	AuthDB        string
	Logger        Logger

	Host string
	Port string
}

func buildConfig(opts ...Option) *Config {
	cfg := &Config{
		NetworkName:   "test-network",
		ContainerName: "mongo-container",
		ImageName:     "mongo:8.0",
		Database:      "test",
		Username:      "root",
		Password:      "root",
		AuthDB:        "admin",
		Logger:        &logger.NoopLogger{},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func defaultHostConfig() func(hc *container.HostConfig) {
	return func(hc *container.HostConfig) {
		hc.AutoRemove = true
	}
}
