package v1

import (
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/service"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer

	inventoryService service.InventoryService
}

func NewAPI(inventoryService service.InventoryService) *api {
	return &api{inventoryService: inventoryService}
}
