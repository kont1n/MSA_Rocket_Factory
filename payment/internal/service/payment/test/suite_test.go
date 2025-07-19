package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/service/payment"
)

type ServiceSuite struct {
	suite.Suite
	ctx     context.Context
	service service.PaymentService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()
	s.service = payment.NewService()
}

func (s *ServiceSuite) TearDownSuite() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
