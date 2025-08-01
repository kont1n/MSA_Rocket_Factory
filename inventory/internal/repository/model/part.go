package model

import (
	"time"
)

type RepositoryPart struct {
	OrderUuid     string           `bson:"order_uuid"`
	Name          string           `bson:"name"`
	Description   string           `bson:"description"`
	Price         float64          `bson:"price"`
	StockQuantity int64            `bson:"stock_quantity"`
	Category      int              `bson:"category"`
	Dimensions    Dimensions       `bson:"dimensions"`
	Manufacturer  Manufacturer     `bson:"manufacturer"`
	Tags          []string         `bson:"tags"`
	Metadata      map[string]Value `bson:"metadata"`
	CreatedAt     time.Time        `bson:"created_at"`
	UpdatedAt     time.Time        `bson:"updated_at"`
}
