package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/service/order"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/service/order/mocks"
)

type ServiceSuite struct {
	suite.Suite
	ctx             context.Context
	service         service.OrderService
	orderRepository *mocks.MockOrderRepository
	inventoryClient *mocks.MockInventoryClient
	paymentClient   *mocks.MockPaymentClient
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()
	s.orderRepository = mocks.NewMockOrderRepository(s.T())
	s.inventoryClient = mocks.NewMockInventoryClient(s.T())
	s.paymentClient = mocks.NewMockPaymentClient(s.T())

	s.service = order.NewService(
		s.orderRepository,
		s.inventoryClient,
		s.paymentClient,
	)
}

func (s *ServiceSuite) SetupTest() {
	// Сбрасываем моки перед каждым тестом
	s.orderRepository.ExpectedCalls = nil
	s.inventoryClient.ExpectedCalls = nil
	s.paymentClient.ExpectedCalls = nil
}

func (s *ServiceSuite) TearDownSuite() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
