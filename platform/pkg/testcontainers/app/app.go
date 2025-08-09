package app

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

const (
	defaultAppName        = "app"
	defaultAppPort        = "50051"
	defaultStartupTimeout = 1 * time.Minute
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Config struct {
	Name          string
	DockerfileDir string
	Dockerfile    string
	Port          string
	Env           map[string]string
	Networks      []string
	LogOutput     io.Writer
	StartupWait   wait.Strategy
	Logger        Logger
}

type Container struct {
	container    testcontainers.Container
	externalHost string
	externalPort string
	cfg          *Config
}

func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	cfg := &Config{
		Name:          defaultAppName,
		Port:          defaultAppPort,
		Dockerfile:    "Dockerfile",
		DockerfileDir: ".",
		LogOutput:     io.Discard,
		StartupWait:   wait.ForListeningPort(defaultAppPort + "/tcp").WithStartupTimeout(defaultStartupTimeout),
		Env:           make(map[string]string),
		Logger:        &logger.NoopLogger{},
	}
	for _, opt := range opts {
		opt(cfg)
	}

	req := testcontainers.ContainerRequest{
		Name: cfg.Name,
		FromDockerfile: testcontainers.FromDockerfile{
			Context:        cfg.DockerfileDir,
			Dockerfile:     cfg.Dockerfile,
			BuildLogWriter: cfg.LogOutput,
		},
		Networks:           cfg.Networks,
		Env:                cfg.Env,
		WaitingFor:         cfg.StartupWait,
		ExposedPorts:       []string{cfg.Port + "/tcp"},
		HostConfigModifier: DefaultHostConfig(),
	}

	genericContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, errors.Errorf("failed to start app genericContainer: %v", err)
	}

	mappedPort, err := genericContainer.MappedPort(ctx, nat.Port(cfg.Port+"/tcp"))
	if err != nil {
		return nil, errors.Errorf("failed to get mapped externalPort: %v", err)
	}

	host, err := genericContainer.Host(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to get genericContainer externalHost: %v", err)
	}

	go streamContainerLogs(ctx, genericContainer, cfg.LogOutput)

	cfg.Logger.Info(ctx, "App container started", zap.String("uri:", net.JoinHostPort(host, mappedPort.Port())))

	return &Container{
		container:    genericContainer,
		externalHost: host,
		externalPort: mappedPort.Port(),
		cfg:          cfg,
	}, nil
}

func (a *Container) Address() string {
	return net.JoinHostPort(a.externalHost, a.externalPort)
}

func (a *Container) Terminate(ctx context.Context) error {
	return a.container.Terminate(ctx)
}

func streamContainerLogs(ctx context.Context, container testcontainers.Container, out io.Writer) {
	logs, err := container.Logs(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get container logs", zap.Error(err))
		return
	}
	defer func() {
		err = logs.Close()
		if err != nil {
			logger.Error(ctx, "failed to close container logs", zap.Error(err))
		}
	}()

	go func() {
		_, err = io.Copy(out, logs)
		if err != nil && !errors.Is(err, io.EOF) {
			logger.Error(ctx, "error copying container logs", zap.Error(err))
		}
	}()
}

func DefaultHostConfig() func(hc *container.HostConfig) {
	return func(hc *container.HostConfig) {
		hc.AutoRemove = true
	}
}
