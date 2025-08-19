package part_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/mocks"
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/service/part"
)

type ServiceSuite struct {
	suite.Suite
	inventoryRepo *mocks.InventoryRepository
	service       service.InventoryService
}

func (s *ServiceSuite) SetupSuite() {
	s.inventoryRepo = mocks.NewInventoryRepository(s.T())
	s.service = part.NewService(s.inventoryRepo)
}

func (s *ServiceSuite) SetupTest() {
	// Сбрасываем моки перед каждым тестом
	s.inventoryRepo.ExpectedCalls = nil
}

func (s *ServiceSuite) TearDownSuite() {
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
