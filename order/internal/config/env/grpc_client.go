package env

import (
	"os"
)

type grpcClientConfig struct {
	inventoryAddress string
	paymentAddress   string
}

func NewGRPCClientConfig() (*grpcClientConfig, error) {
	inventoryAddress := os.Getenv("INVENTORY_GRPC_ADDRESS")
	if inventoryAddress == "" {
		inventoryAddress = "localhost:50051"
	}

	paymentAddress := os.Getenv("PAYMENT_GRPC_ADDRESS")
	if paymentAddress == "" {
		paymentAddress = "localhost:50052"
	}

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
