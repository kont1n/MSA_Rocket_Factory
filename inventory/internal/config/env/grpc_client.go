package env

import (
	"net"
	"os"
)

// grpcClientConfig конфигурация для подключения к IAM сервису
type grpcClientConfig struct {
	iamAddress string // адрес IAM сервиса (например: localhost:50051)
}

// NewIAMConfig создает конфигурацию IAM из переменных окружения
func NewGRPCClientConfig() (*grpcClientConfig, error) {
	host := os.Getenv("IAM_GRPC_HOST")
	if host == "" {
		host = "localhost" // дефолтный хост
	}

	port := os.Getenv("IAM_GRPC_PORT")
	if port == "" {
		port = "50051" // дефолтный порт
	}

	target := net.JoinHostPort(host, port)

	return &grpcClientConfig{
		iamAddress: target,
	}, nil
}

// Target возвращает адрес IAM сервиса
func (c *grpcClientConfig) IAMAddress() string {
	return c.iamAddress
}
