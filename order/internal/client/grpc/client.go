package grpc

import (
	"context"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter *model.Filter) (*[]model.Part, error)
}

type PaymentClient interface {
	CreatePayment(ctx context.Context, order *model.Order) (*model.Order, error)
}
