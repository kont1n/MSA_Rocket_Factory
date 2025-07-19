package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/mocks"
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/service/part"
)

type ServiceSuite struct {
	suite.Suite
	ctx           context.Context
	inventoryRepo *mocks.InventoryRepository
	service       service.InventoryService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()
	s.inventoryRepo = mocks.NewInventoryRepository(s.T())
	s.service = part.NewService(
		s.inventoryRepo,
	)
}

func (s *ServiceSuite) TearDownSuite() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
