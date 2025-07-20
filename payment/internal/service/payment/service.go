package payment

import (
	def "github.com/kont1n/MSA_Rocket_Factory/payment/internal/service"
)

var _ def.PaymentService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
