package app

import (
	"context"

	paymentV1API "github.com/kont1n/MSA_Rocket_Factory/payment/internal/api/payment/v1"
	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/service"
	paymentService "github.com/kont1n/MSA_Rocket_Factory/payment/internal/service/payment"
	paymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentAPIv1   paymentV1.PaymentServiceServer
	paymentService service.PaymentService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1API(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentAPIv1 == nil {
		d.paymentAPIv1 = paymentV1API.NewAPI(d.PaymentService(ctx))
	}
	return d.paymentAPIv1
}

func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService()
	}
	return d.paymentService
}
