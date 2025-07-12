package converter

import (
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

func PartListToModel(parts []*inventoryV1.Part) (*[]model.Part, error) {
	result := make([]*model.Part, 0, len(parts))
	for _, part := range parts {
		result = append(result, PartToModel(part))
	}
	return &result, nil
}

func PartToModel(part *inventoryV1.Part) *model.Part {
	return &model.Part{
		PartUUID: part.PartUuid,
		PartName: part.Name,
		PartDescription: part.Description,
		PartPrice: part.Price,
		PartQuantity: part.StockQuantity,
		PartCategory: part.Category,
		PartManufacturer: part.Manufacturer.Name,
		PartManufacturerCountry: part.Manufacturer.Country,
	}
}