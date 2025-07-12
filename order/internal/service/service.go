package service

import (
	"context"
	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, error)
	PayOrder(ctx context.Context, order *model.Order) (*model.Order, error)
	CancelOrder(ctx context.Context, order *model.Order) (*model.Order, error)
}
