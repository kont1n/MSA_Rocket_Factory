package inmemory

import (
	"sync"

	def "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]*repoModel.RepositoryPart
}

func NewRepository() *repository {
	repo := &repository{
		data: make(map[string]*repoModel.RepositoryPart),
	}

	repo.addTestData()

	return repo
}
