package env

import (
	"net"
	"os"
)

type grpcClientConfig struct {
	iamAddress string
}

func NewGRPCClientConfig() (*grpcClientConfig, error) {
	iamHost := os.Getenv("IAM_GRPC_HOST")
	if iamHost == "" {
		iamHost = "localhost"
	}

	iamPort := os.Getenv("IAM_GRPC_PORT")
	if iamPort == "" {
		iamPort = "50051"
	}

	iamAddress := net.JoinHostPort(iamHost, iamPort)

	return &grpcClientConfig{
		iamAddress: iamAddress,
	}, nil
}

func (c *grpcClientConfig) IAMAddress() string {
	return c.iamAddress
}
