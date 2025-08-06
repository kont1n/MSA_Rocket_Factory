package converter

import (
	"github.com/google/uuid"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

func ToModelPartsList(parts []*inventoryV1.Part) (*[]model.Part, error) {
	result := make([]model.Part, 0, len(parts))
	for _, part := range parts {
		modelPart, err := ToModelPart(part)
		if err != nil {
			return nil, err
		}
		result = append(result, *modelPart)
	}
	return &result, nil
}

func ToModelPart(part *inventoryV1.Part) (*model.Part, error) {
	id, err := uuid.Parse(part.PartUuid)
	if err != nil {
		return nil, model.ErrConvertFromClient
	}
	return &model.Part{
		PartUUID:    id,
		Name:        part.Name,
		Description: part.Description,
		Price:       part.Price,
	}, nil
}
