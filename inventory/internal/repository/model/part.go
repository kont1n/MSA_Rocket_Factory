package model

import (
	"time"
)

type RepositoryPart struct {
	OrderUuid     string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      int
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          []string
	Metadata      map[string]Value
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
