package inmemory_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/repository"
	inmemoryRepo "github.com/kont1n/MSA_Rocket_Factory/order/internal/repository/inmemory"
)

type InMemoryOrderRepositorySuite struct {
	suite.Suite
	repository repository.OrderRepository
}

func (s *InMemoryOrderRepositorySuite) SetupTest() {
	// Создаем новый репозиторий для каждого теста
	s.repository = inmemoryRepo.NewRepository()
}

func (s *InMemoryOrderRepositorySuite) TearDownTest() {
	// Очистка не требуется, так как каждый тест создает новый репозиторий
}

func TestInMemoryOrderRepository(t *testing.T) {
	suite.Run(t, new(InMemoryOrderRepositorySuite))
}
