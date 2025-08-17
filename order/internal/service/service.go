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
	UpdateOrderStatus(ctx context.Context, orderUUID string, status model.OrderStatus) error
}

type OrderPaidProducer interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error
}

type ShipAssembledConsumer interface {
	RunConsumer(ctx context.Context) error
}
