package mongo

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func startMongoContainer(ctx context.Context, cfg *Config) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Name:     cfg.ContainerName,
		Image:    cfg.ImageName,
		Networks: []string{cfg.NetworkName},
		Env: map[string]string{
			mongoEnvUsernameKey: cfg.Username,
			mongoEnvPasswordKey: cfg.Password,
		},
		WaitingFor:         wait.ForListeningPort(mongoPort + "/tcp").WithStartupTimeout(mongoStartupTimeout),
		HostConfigModifier: defaultHostConfig(),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, errors.Errorf("failed to start mongo container: %v", err)
	}

	return container, nil
}

func getContainerHostPort(ctx context.Context, container testcontainers.Container) (string, string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", "", errors.Errorf("failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, mongoPort+"/tcp")
	if err != nil {
		return "", "", errors.Errorf("failed to get mapped port: %v", err)
	}

	return host, port.Port(), nil
}

func buildMongoURI(cfg *Config) string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.AuthDB,
	)
}
