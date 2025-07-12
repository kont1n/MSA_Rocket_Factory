package inmemory

import (
	"sync"

	def "github.com/kont1n/MSA_Rocket_Factory/order/internal/repository"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/order/internal/repository/model"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.Order
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]repoModel.Order),
	}
}
