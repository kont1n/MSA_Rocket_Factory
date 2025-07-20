package converter

import (
	"strconv"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

func PartsFilterToProto(filter *model.Filter) *inventoryV1.PartsFilter {
	uuids := make([]string, 0, len(filter.PartUUIDs))
	for _, id := range filter.PartUUIDs {
		uuids = append(uuids, id.String())
	}

	categories := make([]inventoryV1.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		// Безопасная конвертация int в int32 через строку
		categoryStr := strconv.Itoa(int(category))
		if categoryInt32, err := strconv.ParseInt(categoryStr, 10, 32); err == nil {
			categories = append(categories, inventoryV1.Category(int32(categoryInt32)))
		}
	}

	return &inventoryV1.PartsFilter{
		PartUuid:            uuids,
		PartName:            filter.PartNames,
		Category:            categories,
		ManufacturerCountry: filter.ManufacturerCountries,
		Tags:                filter.Tags,
	}
}
