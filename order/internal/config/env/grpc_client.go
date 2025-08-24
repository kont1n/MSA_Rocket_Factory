package env

import (
	"net"
	"os"
)

type grpcClientConfig struct {
	inventoryAddress string
	paymentAddress   string
	iamAddress       string
}

func NewGRPCClientConfig() (*grpcClientConfig, error) {
	inventoryHost := os.Getenv("INVENTORY_GRPC_HOST")
	if inventoryHost == "" {
		inventoryHost = "localhost"
	}

	inventoryPort := os.Getenv("INVENTORY_GRPC_PORT")
	if inventoryPort == "" {
		inventoryPort = "50051"
	}

	inventoryAddress := net.JoinHostPort(inventoryHost, inventoryPort)

	paymentHost := os.Getenv("PAYMENT_GRPC_HOST")
	if paymentHost == "" {
		paymentHost = "localhost"
	}

	paymentPort := os.Getenv("PAYMENT_GRPC_PORT")
	if paymentPort == "" {
		paymentPort = "50052"
	}

	paymentAddress := net.JoinHostPort(paymentHost, paymentPort)

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
		inventoryAddress: inventoryAddress,
		paymentAddress:   paymentAddress,
		iamAddress:       iamAddress,
	}, nil
}

func (c *grpcClientConfig) InventoryAddress() string {
	return c.inventoryAddress
}

func (c *grpcClientConfig) PaymentAddress() string {
	return c.paymentAddress
}

func (c *grpcClientConfig) IAMAddress() string {
	return c.iamAddress
}
