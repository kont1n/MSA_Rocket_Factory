package model

import "github.com/google/uuid"

type Filter struct {
	PartUUIDs             []uuid.UUID
	PartNames             []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
