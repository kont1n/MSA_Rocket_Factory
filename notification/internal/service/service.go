package service

import (
	"context"

	"github.com/kont1n/MSA_Rocket_Factory/notification/internal/model"
)

type NotificationService interface {
	NotifyOrderPaid(ctx context.Context, event *model.OrderPaidEvent) error
	NotifyShipAssembled(ctx context.Context, event *model.ShipAssembledEvent) error
}

type OrderPaidConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type ShipAssembledConsumerService interface {
	RunConsumer(ctx context.Context) error
}
