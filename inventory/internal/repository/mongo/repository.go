package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	db *mongo.Database
}

func NewRepository(database *mongo.Database) *repository {
	return &repository{
		db: database,
	}
}
