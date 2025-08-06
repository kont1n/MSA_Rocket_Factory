package converter

import (
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProtoPart(part *model.Part) *inventoryV1.Part {
	// Конвертируем Category из model.Category в inventoryV1.Category
	var protoCategory inventoryV1.Category
	switch part.Category {
	case model.ENGINE:
		protoCategory = inventoryV1.Category_CATEGORY_ENGINE
	case model.FUEL:
		protoCategory = inventoryV1.Category_CATEGORY_FUEL
	case model.PORTHOLE:
		protoCategory = inventoryV1.Category_CATEGORY_PORTHOLE
	case model.WING:
		protoCategory = inventoryV1.Category_CATEGORY_WING
	default:
		protoCategory = inventoryV1.Category_CATEGORY_UNSPECIFIED
	}

	// Конвертируем Dimensions
	dimensions := &inventoryV1.Dimensions{
		Length: part.Dimensions.Length,
		Width:  part.Dimensions.Width,
		Height: part.Dimensions.Height,
		Weight: part.Dimensions.Weight,
	}

	// Конвертируем Manufacturer
	manufacturer := &inventoryV1.Manufacturer{
		Name:    part.Manufacturer.Name,
		Country: part.Manufacturer.Country,
		Url:     part.Manufacturer.Website,
	}

	// Конвертируем Metadata
	metadata := make(map[string]*inventoryV1.Value)
	for key, value := range part.Metadata {
		protoValue := &inventoryV1.Value{}
		switch {
		case value.StringValue != "":
			protoValue.Kind = &inventoryV1.Value_StringValue{StringValue: value.StringValue}
		case value.Int64Value != 0:
			protoValue.Kind = &inventoryV1.Value_Int64Value{Int64Value: value.Int64Value}
		case value.Float64Value != 0:
			protoValue.Kind = &inventoryV1.Value_DoubleValue{DoubleValue: value.Float64Value}
		case value.BoolValue:
			protoValue.Kind = &inventoryV1.Value_BoolValue{BoolValue: value.BoolValue}
		}
		metadata[key] = protoValue
	}

	return &inventoryV1.Part{
		PartUuid:      part.OrderUuid.String(),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      protoCategory,
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          part.Tags,
		Metadata:      metadata,
		CreatedAt:     timestamppb.New(part.CreatedAt),
		UpdatedAt:     timestamppb.New(part.UpdatedAt),
	}
}
