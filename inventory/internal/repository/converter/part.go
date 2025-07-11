package converter

import (
	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/model"
)

func RepoToModel(repoPart *repoModel.RepositoryPart) (part *model.Part) {
	uuid, err := uuid.Parse(repoPart.OrderUuid)
	if err != nil {
	}

	manufacturer := model.Manufacturer{
		Name:    repoPart.Manufacturer.Name,
		Country: repoPart.Manufacturer.Country,
		Website: repoPart.Manufacturer.Website,
	}

	dimension := model.Dimensions{
		Length: repoPart.Dimensions.Length,
		Width:  repoPart.Dimensions.Width,
		Height: repoPart.Dimensions.Height,
		Weight: repoPart.Dimensions.Weight,
	}

	metadata := make(map[string]model.Value)
	for key, value := range repoPart.Metadata {
		metadata[key] = model.Value{
			StringValue:  value.StringValue,
			Int64Value:   value.Int64Value,
			Float64Value: value.Float64Value,
			BoolValue:    value.BoolValue,
		}
	}

	return &model.Part{
		OrderUuid:     uuid,
		Name:          repoPart.Name,
		Description:   repoPart.Description,
		Price:         repoPart.Price,
		StockQuantity: repoPart.StockQuantity,
		Category:      model.ToCategory(repoPart.Category),
		Dimensions:    dimension,
		Manufacturer:  manufacturer,
		Tags:          repoPart.Tags,
		Metadata:      metadata,
		CreatedAt:     repoPart.CreatedAt,
		UpdatedAt:     repoPart.UpdatedAt,
	}
}
