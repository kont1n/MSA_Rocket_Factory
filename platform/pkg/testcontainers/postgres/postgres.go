package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

// Container представляет PostgreSQL контейнер
type Container struct {
	container testcontainers.Container
	config    *Config
}

// NewContainer создает и запускает новый PostgreSQL контейнер
func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	config := NewConfig()

	// Применяем опции
	for _, opt := range opts {
		opt(config)
	}

	// Настройки контейнера
	req := testcontainers.ContainerRequest{
		Image:        config.ImageName,
		ExposedPorts: []string{config.Port + "/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       config.Database,
			"POSTGRES_USER":     config.Username,
			"POSTGRES_PASSWORD": config.Password,
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(nat.Port(config.Port+"/tcp")),
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	}

	// Добавляем имя контейнера если указано
	if config.ContainerName != "" {
		req.Name = config.ContainerName
	}

	// Добавляем сеть если указана
	if config.NetworkName != "" {
		req.Networks = []string{config.NetworkName}
	}

	// Создаем и запускаем контейнер
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось создать PostgreSQL контейнер: %w", err)
	}

	config.Logger.Info(ctx, "PostgreSQL контейнер успешно запущен",
		zap.String("container_name", config.ContainerName),
		zap.String("database", config.Database))

	return &Container{
		container: container,
		config:    config,
	}, nil
}

// Config возвращает конфигурацию контейнера
func (c *Container) Config() *Config {
	return c.config
}

// Host возвращает хост для подключения к PostgreSQL
func (c *Container) Host(ctx context.Context) (string, error) {
	return c.container.Host(ctx)
}

// Port возвращает порт для подключения к PostgreSQL
func (c *Container) Port(ctx context.Context) (nat.Port, error) {
	return c.container.MappedPort(ctx, nat.Port(c.config.Port+"/tcp"))
}

// ConnectionString возвращает строку подключения к PostgreSQL
func (c *Container) ConnectionString(ctx context.Context) (string, error) {
	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}

	port, err := c.Port(ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.config.Username,
		c.config.Password,
		host,
		port.Port(),
		c.config.Database,
	), nil
}

// Terminate останавливает и удаляет контейнер
func (c *Container) Terminate(ctx context.Context) error {
	c.config.Logger.Info(ctx, "Остановка PostgreSQL контейнера",
		zap.String("container_name", c.config.ContainerName))

	if err := c.container.Terminate(ctx); err != nil {
		c.config.Logger.Error(ctx, "не удалось остановить PostgreSQL контейнер", zap.Error(err))
		return fmt.Errorf("не удалось остановить PostgreSQL контейнер: %w", err)
	}

	return nil
}
