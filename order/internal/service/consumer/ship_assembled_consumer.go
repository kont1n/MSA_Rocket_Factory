package consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	serviceDef "github.com/kont1n/MSA_Rocket_Factory/order/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/kafka"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

type service struct {
	shipAssembledConsumer kafka.Consumer
	shipAssembledDecoder  ShipAssembledDecoder
	orderService          serviceDef.OrderService
	metrics               KafkaMetrics
}

// KafkaMetrics интерфейс для метрик Kafka
type KafkaMetrics interface {
	RecordConsumerMessage(ctx context.Context, topic string, partition int32, groupID string, success bool)
	RecordConsumerLag(ctx context.Context, topic string, partition int32, groupID string, lagSeconds float64, offsetLag int64)
}

type ShipAssembledDecoder interface {
	Decode(data []byte) (*model.ShipAssembledEvent, error)
}

func NewService(shipAssembledConsumer kafka.Consumer, shipAssembledDecoder ShipAssembledDecoder, orderService serviceDef.OrderService, metrics KafkaMetrics) *service {
	return &service{
		shipAssembledConsumer: shipAssembledConsumer,
		shipAssembledDecoder:  shipAssembledDecoder,
		orderService:          orderService,
		metrics:               metrics,
	}
}

// SetOrderService устанавливает OrderService для отправки событий
func (s *service) SetOrderService(orderService serviceDef.OrderService) {
	s.orderService = orderService
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
	// Записываем метрики при начале обработки
	if s.metrics != nil {
		s.metrics.RecordConsumerMessage(ctx, "ship-assembled", 0, "order-ship-assembled-group", true)
	}

	event, err := s.shipAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ShipAssembled", zap.Error(err))
		// Записываем ошибку в метрики
		if s.metrics != nil {
			s.metrics.RecordConsumerMessage(ctx, "ship-assembled", 0, "order-ship-assembled-group", false)
		}
		return err
	}

	logger.Info(ctx, "Processing ShipAssembled message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID.String()),
		zap.String("order_uuid", event.OrderUUID.String()),
	)

	// Обновляем статус заказа на ASSEMBLED (если orderService доступен)
	if s.orderService != nil {
		err = s.orderService.UpdateOrderStatus(ctx, event.OrderUUID.String(), model.StatusAssembled)
		if err != nil {
			logger.Error(ctx, "Failed to update order status to ASSEMBLED", zap.Error(err))
			// Записываем ошибку в метрики
			if s.metrics != nil {
				s.metrics.RecordConsumerMessage(ctx, "ship-assembled", 0, "order-ship-assembled-group", false)
			}
			return err
		}
	} else {
		logger.Warn(ctx, "OrderService is not available, skipping order status update",
			zap.String("order_uuid", event.OrderUUID.String()))
	}

	logger.Info(ctx, "Order status updated to ASSEMBLED successfully",
		zap.String("order_uuid", event.OrderUUID.String()))

	return nil
}
