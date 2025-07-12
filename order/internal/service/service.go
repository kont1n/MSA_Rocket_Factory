package service

import (
	"context"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error)
}