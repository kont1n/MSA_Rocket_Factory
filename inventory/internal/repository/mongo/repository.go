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

func NewRepository(ctx context.Context, database *mongo.Database) def.InventoryRepository {
	repo := &repository{
		db: database,
	}

	// Добавляем тестовые данные при инициализации
	if err := repo.AddTestData(ctx); err != nil {
		log.Printf("Предупреждение: не удалось добавить тестовые данные в MongoDB: %v", err)
	}

	return repo
}
