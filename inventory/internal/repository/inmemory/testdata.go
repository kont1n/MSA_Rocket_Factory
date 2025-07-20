package inmemory

import (
	"log"
	"time"

	repoModel "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/model"
)

// TestData Добавление тестовых данных
func TestData(repo *repository) {
	log.Printf("Add Test Data for inventory service")

	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.data = map[string]*repoModel.RepositoryPart{
		"d973e963-b7e6-4323-8f4e-4bfd5ab8e834": {
			OrderUuid:     "d973e963-b7e6-4323-8f4e-4bfd5ab8e834",
			Name:          "Detail 1",
			Description:   "Detail 1 description",
			Price:         100,
			StockQuantity: 10,
			Category:      1, // ENGINE
			Dimensions: repoModel.Dimensions{
				Length: 100,
				Width:  100,
				Height: 100,
				Weight: 100,
			},
			Manufacturer: repoModel.Manufacturer{
				Country: "China",
				Name:    "Details Fabric",
			},
			Tags:      []string{"tag1", "tag2"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		"d973e963-b7e6-4323-8f4e-4bfd5ab8e835": {
			OrderUuid:     "d973e963-b7e6-4323-8f4e-4bfd5ab8e835",
			Name:          "Detail 2",
			Description:   "Detail 2 description",
			Price:         200,
			StockQuantity: 20,
			Category:      1, // ENGINE
			Dimensions: repoModel.Dimensions{
				Length: 100,
				Width:  100,
				Height: 100,
				Weight: 100,
			},
			Manufacturer: repoModel.Manufacturer{
				Country: "USA",
				Name:    "Details Fabric",
			},
			Tags:      []string{"tag1", "tag2"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}
