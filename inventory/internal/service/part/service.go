package part

import (
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository"
	def "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	repo repository.InventoryRepository
}

func NewService(repo repository.InventoryRepository) *service {
	return &service{
		repo: repo,
	}
}
