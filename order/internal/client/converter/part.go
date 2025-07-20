package converter

import (
	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

func PartListToModel(parts []*inventoryV1.Part) (*[]model.Part, error) {
	result := make([]model.Part, 0, len(parts))
	for _, part := range parts {
		modelPart, err := PartToModel(part)
		if err != nil {
			return nil, err
		}
		result = append(result, *modelPart)
	}
	return &result, nil
}

func PartToModel(part *inventoryV1.Part) (*model.Part, error) {
	id, err := uuid.Parse(part.PartUuid)
	if err != nil {
		return nil, err
	}
	return &model.Part{
		PartUUID:    id,
		Name:        part.Name,
		Description: part.Description,
		Price:       part.Price,
	}, nil
}
