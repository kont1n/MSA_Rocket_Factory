package model

import "github.com/google/uuid"

type Filter struct {
	Uuids                 []uuid.UUID
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
