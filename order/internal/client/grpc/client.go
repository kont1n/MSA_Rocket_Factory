package grpc

import (
	"context"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

type InventoryClient interface {
	GetPart(ctx context.Context, uuid uuid.UUID) (*model.Part, error)
	ListParts(ctx context.Context, filter *model.Filter) (*[]model.Part, error)
}

type PaymentClient interface {
	CreatePayment(ctx context.Context, order *model.Order) (*model.Order, error)
}