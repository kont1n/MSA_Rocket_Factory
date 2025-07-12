package v1

import (
	generaredInventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

type inventoryClient struct {
	generatedClient generaredInventoryV1.InventoryServiceClient
}

func NewClient(generatedClient generaredInventoryV1.InventoryServiceClient) *inventoryClient {
	return &inventoryClient{
		generatedClient: generatedClient,
	}
}
