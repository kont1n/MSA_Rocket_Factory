package env

import (
	"fmt"
	"os"
)

// iamConfig конфигурация для подключения к IAM сервису
type iamConfig struct {
	iamAddress string // адрес IAM сервиса (например: localhost:50051)
}

// NewIAMConfig создает конфигурацию IAM из переменных окружения
func NewIAMConfig() (*iamConfig, error) {
	host := os.Getenv("IAM_GRPC_HOST")
	if host == "" {
		host = "localhost" // дефолтный хост
	}

	port := os.Getenv("IAM_GRPC_PORT")
	if port == "" {
		port = "50051" // дефолтный порт
	}

	target := fmt.Sprintf("%s:%s", host, port)

	return &iamConfig{
		iamAddress: target,
	}, nil
}

// Target возвращает адрес IAM сервиса
func (c *iamConfig) IAMAddress() string {
	return c.iamAddress
}

// String возвращает строковое представление конфигурации
func (c *iamConfig) String() string {
	return fmt.Sprintf("IAM Service Target: %s", c.iamAddress)
}
