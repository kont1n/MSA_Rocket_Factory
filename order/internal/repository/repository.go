package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) (*model.Order, error)
	DeleteOrder(ctx context.Context, id uuid.UUID) error
}