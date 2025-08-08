package app

import (
	"context"
	"fmt"

	inventoryRepository "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/mongo"
	inventoryService "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/service/part"
	paymentV1API "github.com/kont1n/MSA_Rocket_Factory/payment/internal/api/v1"
	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/config"
	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/service"
	paymentService "github.com/kont1n/MSA_Rocket_Factory/payment/internal/service/payment"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/closer"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type diContainer struct {
	paymentAPIv1   paymentV1.PaymentServiceServer
	paymentService service.PaymentService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1API(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.inventoryAPIv1 == nil {
		d.inventoryAPIv1 = paymentV1API.NewAPI(d.PartService(ctx))
	}
	return d.inventoryAPIv1
}

func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService()
	}
	return d.paymentService
}
