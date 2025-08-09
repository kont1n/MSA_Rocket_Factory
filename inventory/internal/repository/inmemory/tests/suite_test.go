package inmemory_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository"
	inmemoryRepo "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/inmemory"
)

type InMemoryRepositorySuite struct {
	suite.Suite
	repository repository.InventoryRepository
}

func (s *InMemoryRepositorySuite) SetupTest() {
	// Создаем новый репозиторий для каждого теста
	s.repository = inmemoryRepo.NewRepository()
}

func (s *InMemoryRepositorySuite) TearDownTest() {
	// Очистка не требуется, так как каждый тест создает новый репозиторий
}

func TestInMemoryRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(InMemoryRepositorySuite))
}
