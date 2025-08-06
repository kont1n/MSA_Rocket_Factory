package converter

import (
	"github.com/google/uuid"
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

func ToModelPart(req *inventoryV1.ListPartsRequest) *model.Filter {
	if req.Filter != nil {
		uuids := make([]uuid.UUID, 0, len(req.Filter.PartUuid))
		for _, uuidStr := range req.Filter.PartUuid {
			if uuid, err := uuid.Parse(uuidStr); err == nil {
				uuids = append(uuids, uuid)
			}
		}

		categories := make([]model.Category, 0, len(req.Filter.Category))
		for _, protoCategory := range req.Filter.Category {
			switch protoCategory {
			case inventoryV1.Category_CATEGORY_ENGINE:
				categories = append(categories, model.ENGINE)
			case inventoryV1.Category_CATEGORY_FUEL:
				categories = append(categories, model.FUEL)
			case inventoryV1.Category_CATEGORY_PORTHOLE:
				categories = append(categories, model.PORTHOLE)
			case inventoryV1.Category_CATEGORY_WING:
				categories = append(categories, model.WING)
			}
		}

		filter := &model.Filter{
			Uuids:                 uuids,
			Names:                 req.Filter.PartName,
			Categories:            categories,
			ManufacturerCountries: req.Filter.ManufacturerCountry,
			Tags:                  req.Filter.Tags,
		}
		return filter
	}

	return nil
}
