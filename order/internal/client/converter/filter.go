package converter

import (
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

func PartsFilterToProto(filter *model.Filter) *inventoryV1.Filter {
	categories := make([]inventoryV1.Category, 0, len(filter.PartCategories))
	for _, category := range filter.PartCategories {
		categories = append(categories, inventoryV1.Category(category))
	}

	return &inventoryV1.PartsFilter{
		PartUuids:           filter.PartUUIDs,
		PartNames:           filter.PartNames,
		Category:            categories,
		ManufacturerCountry: filter.PartManufacturerCountry,
		Tag:                 filter.Tags,
	}
}
