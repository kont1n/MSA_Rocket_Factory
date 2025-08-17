package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/client/grpc"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/repository"
	def "github.com/kont1n/MSA_Rocket_Factory/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository   repository.OrderRepository
	inventoryClient   grpc.InventoryClient
	paymentClient     grpc.PaymentClient
	orderPaidProducer def.OrderPaidProducer
}

func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,
	orderPaidProducer def.OrderPaidProducer,
) *service {
	return &service{
		orderRepository:   orderRepository,
		inventoryClient:   inventoryClient,
		paymentClient:     paymentClient,
		orderPaidProducer: orderPaidProducer,
	}
}

func (s *service) UpdateOrderStatus(ctx context.Context, orderUUID string, status model.OrderStatus) error {
	uuid, err := uuid.Parse(orderUUID)
	if err != nil {
		return fmt.Errorf("invalid order UUID: %w", err)
	}

	order, err := s.orderRepository.GetOrder(ctx, uuid)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	order.Status = status

	_, err = s.orderRepository.UpdateOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}
