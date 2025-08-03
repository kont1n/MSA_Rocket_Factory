package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository"
)

var _ def.InventoryRepository = (*repository)(nil)

const partsCollection = "parts"

type repository struct {
	db *mongo.Database
}

func NewRepository(database *mongo.Database) *repository {
	repo := &repository{
		db: database,
	}

	// Добавляем тестовые данные при инициализации
	if err := repo.AddTestData(context.Background()); err != nil {
		log.Printf("Предупреждение: не удалось добавить тестовые данные в MongoDB: %v", err)
	}

	return repo
}
