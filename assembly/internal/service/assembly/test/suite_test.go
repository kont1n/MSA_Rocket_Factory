package assembly_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/assembly/internal/service/assembly"
)

type AssemblyServiceSuite struct {
	suite.Suite
	assemblyProducerService *MockProducerService
	service                 service.AssemblyService
}

type MockProducerService struct {
	ProduceAssemblyFunc func(ctx context.Context, event model.ShipAssembledEvent) error
}

func (m *MockProducerService) ProduceAssembly(ctx context.Context, event model.ShipAssembledEvent) error {
	if m.ProduceAssemblyFunc != nil {
		return m.ProduceAssemblyFunc(ctx, event)
	}
	return nil
}

func (s *AssemblyServiceSuite) SetupSuite() {
	s.assemblyProducerService = &MockProducerService{}
	s.service = assembly.NewService(s.assemblyProducerService)
}

func (s *AssemblyServiceSuite) SetupTest() {
	// Сбрасываем моки перед каждым тестом
	s.assemblyProducerService.ProduceAssemblyFunc = nil
}

func (s *AssemblyServiceSuite) TearDownSuite() {
}

func TestAssemblyService(t *testing.T) {
	suite.Run(t, new(AssemblyServiceSuite))
}
