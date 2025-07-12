package v1

import (
	generaredInventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
	def "github.com/kont1n/MSA_Rocket_Factory/order/internal/client"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	generatedClient generaredInventoryV1.InventoryServiceClient
}

func NewClient(generatedClient generaredInventoryV1.InventoryServiceClient) def.InventoryClient {
	return &client{
		generatedClient: generatedClient,
	}
}
