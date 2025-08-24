package consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type service struct {
	shipAssembledConsumer kafka.Consumer
	shipAssembledDecoder  ShipAssembledDecoder
	orderService          OrderService
}

type ShipAssembledDecoder interface {
	Decode(data []byte) (*model.ShipAssembledEvent, error)
}

type OrderService interface {
	UpdateOrderStatus(ctx context.Context, orderUUID string, status model.OrderStatus) error
}

func NewService(shipAssembledConsumer kafka.Consumer, shipAssembledDecoder ShipAssembledDecoder, orderService OrderService) *service {
	return &service{
		shipAssembledConsumer: shipAssembledConsumer,
		shipAssembledDecoder:  shipAssembledDecoder,
		orderService:          orderService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting ShipAssembled Consumer service")

	err := s.shipAssembledConsumer.Consume(ctx, s.ShipAssembledHandler)
	if err != nil {
		logger.Error(ctx, "Consume from ship.assembled topic error", zap.Error(err))
		return err
	}

	return nil
}

func (s *service) ShipAssembledHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.shipAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ShipAssembled", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing ShipAssembled message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID.String()),
		zap.String("order_uuid", event.OrderUUID.String()),
	)

	// Обновляем статус заказа на ASSEMBLED
	err = s.orderService.UpdateOrderStatus(ctx, event.OrderUUID.String(), model.StatusAssembled)
	if err != nil {
		logger.Error(ctx, "Failed to update order status to ASSEMBLED", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Order status updated to ASSEMBLED successfully",
		zap.String("order_uuid", event.OrderUUID.String()))

	return nil
}
