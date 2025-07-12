package model

import (
	"github.com/google/uuid"
)

type Part struct {
	PartUUID    uuid.UUID
	Name        string
	Description string
	Price       float64
}
