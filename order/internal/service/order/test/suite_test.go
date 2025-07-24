package order_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	clientMocks "github.com/kont1n/MSA_Rocket_Factory/order/internal/client/grpc/mocks"
	repoMocks "github.com/kont1n/MSA_Rocket_Factory/order/internal/repository/mocks"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/service/order"
)

type ServiceSuite struct {
	suite.Suite
	service         service.OrderService
	orderRepository *repoMocks.OrderRepository
	inventoryClient *clientMocks.InventoryClient
	paymentClient   *clientMocks.PaymentClient
}

func (s *ServiceSuite) SetupSuite() {
	s.orderRepository = repoMocks.NewOrderRepository(s.T())
	s.inventoryClient = clientMocks.NewInventoryClient(s.T())
	s.paymentClient = clientMocks.NewPaymentClient(s.T())

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
