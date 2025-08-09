package payment_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/service/payment"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type ServiceSuite struct {
	suite.Suite
	service service.PaymentService
}

func (s *ServiceSuite) SetupSuite() {
	// Инициализируем логгер для тестов
	err := logger.Init("info", false)
	if err != nil {
		s.T().Fatalf("Failed to initialize logger: %v", err)
	}

	s.service = payment.NewService()
}

func (s *ServiceSuite) TearDownSuite() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
