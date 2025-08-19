package order_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	clientMocks "github.com/kont1n/MSA_Rocket_Factory/order/internal/client/grpc/mocks"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	repoMocks "github.com/kont1n/MSA_Rocket_Factory/order/internal/repository/mocks"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/service/order"
)

type ServiceSuite struct {
	suite.Suite
	service           service.OrderService
	orderRepository   *repoMocks.OrderRepository
	inventoryClient   *clientMocks.InventoryClient
	paymentClient     *clientMocks.PaymentClient
	orderPaidProducer *mockOrderPaidProducer
}

func (s *ServiceSuite) SetupSuite() {
	s.orderRepository = repoMocks.NewOrderRepository(s.T())
	s.inventoryClient = clientMocks.NewInventoryClient(s.T())
	s.paymentClient = clientMocks.NewPaymentClient(s.T())

	// Создаем мок для OrderPaidProducer
	s.orderPaidProducer = &mockOrderPaidProducer{}

	s.service = order.NewService(
		s.orderRepository,
		s.inventoryClient,
		s.paymentClient,
		s.orderPaidProducer,
	)
}

func (s *ServiceSuite) SetupTest() {
	// Сбрасываем моки перед каждым тестом
	s.orderRepository.ExpectedCalls = nil
	s.inventoryClient.ExpectedCalls = nil
	s.paymentClient.ExpectedCalls = nil
	s.orderPaidProducer.lastEvent = nil
}

func (s *ServiceSuite) TearDownSuite() {
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

// mockOrderPaidProducer - мок для OrderPaidProducer
type mockOrderPaidProducer struct {
	lastEvent *model.OrderPaidEvent
}

func (m *mockOrderPaidProducer) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	m.lastEvent = &event
	return nil
}

func (m *mockOrderPaidProducer) GetLastEvent() *model.OrderPaidEvent {
	return m.lastEvent
}
