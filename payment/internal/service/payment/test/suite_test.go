package payment_test

import (
	"testing"

	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/service/payment"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	service service.PaymentService
}

func (s *ServiceSuite) SetupSuite() {
	s.service = payment.NewService()
}

func (s *ServiceSuite) TearDownSuite() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
