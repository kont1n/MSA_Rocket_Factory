package consumer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type ConsumerServiceSuite struct {
	suite.Suite
}

func (s *ConsumerServiceSuite) SetupSuite() {
	// Инициализируем no-op логгер для тестов
	logger.SetNopLogger()
}

// mockNotificationService - мок для NotificationService
type mockNotificationService struct {
	mock.Mock
}

func (m *mockNotificationService) NotifyOrderPaid(ctx context.Context, event *model.OrderPaidEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockNotificationService) NotifyShipAssembled(ctx context.Context, event *model.ShipAssembledEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockNotificationService) reset() {
	m.ExpectedCalls = nil
}

func (s *ConsumerServiceSuite) SetupTest() {
	// Сбрасываем моки перед каждым тестом
}

func (s *ConsumerServiceSuite) TearDownSuite() {
}

func TestConsumerServiceSuite(t *testing.T) {
	suite.Run(t, new(ConsumerServiceSuite))
}
