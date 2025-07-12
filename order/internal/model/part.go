package model

import "github.com/google/uuid"

type Part struct {
	PartUUID uuid.UUID
	PartName string
	PartDescription string
	PartPrice float32
	PartQuantity int32
	PartCategory string
	PartManufacturer string
	PartManufacturerCountry string
}	