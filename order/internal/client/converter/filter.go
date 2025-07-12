package converter

import (
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

func PartsFilterToProto(filter model.Filter) *inventoryV1.PartsFilter {
	uuids := make([]string, 0, len(filter.PartUUIDs))
	for _, id := range filter.PartUUIDs {
		uuids = append(uuids, id.String())
	}

	categories := make([]inventoryV1.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(categories, inventoryV1.Category(category))
	}

	return &inventoryV1.PartsFilter{
		PartUuid:            uuids,
		PartName:            filter.PartNames,
		Category:            categories,
		ManufacturerCountry: filter.ManufacturerCountries,
		Tags:                filter.Tags,
	}
}
