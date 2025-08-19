package env

import (
	"os"
)

type grpcClientConfig struct {
	inventoryAddress string
	paymentAddress   string
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

	inventoryAddress := inventoryHost + ":" + inventoryPort

	paymentHost := os.Getenv("PAYMENT_GRPC_HOST")
	if paymentHost == "" {
		paymentHost = "localhost"
	}

	paymentPort := os.Getenv("PAYMENT_GRPC_PORT")
	if paymentPort == "" {
		paymentPort = "50052"
	}

	paymentAddress := paymentHost + ":" + paymentPort

	return &grpcClientConfig{
		inventoryAddress: inventoryAddress,
		paymentAddress:   paymentAddress,
	}, nil
}

func (c *grpcClientConfig) InventoryAddress() string {
	return c.inventoryAddress
}

func (c *grpcClientConfig) PaymentAddress() string {
	return c.paymentAddress
}
