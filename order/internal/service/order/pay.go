package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
)

func (s service) PayOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	// Получаем заказ по UUID
	dbOrder, err := s.orderRepository.GetOrder(ctx, order.OrderUUID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get order from repository: %w", err)
	}

	// Устанавливаем метод оплаты из запроса
	dbOrder.PaymentMethod = order.PaymentMethod

	// Выполняем запрос к API для оплаты заказа
	paidOrder, err := s.paymentClient.CreatePayment(ctx, dbOrder)
	if err != nil {
		return nil, fmt.Errorf("service: failed to create payment in payment client: %w", err)
	}

	// Обновляем заказ в хранилище
	updatedOrder, err := s.orderRepository.UpdateOrder(ctx, paidOrder)
	if err != nil {
		return nil, fmt.Errorf("service: failed to update order in repository: %w", err)
	}

	// Отправляем событие OrderPaid
	event := model.OrderPaidEvent{
		EventUUID:       uuid.New(),
		OrderUUID:       updatedOrder.OrderUUID,
		UserUUID:        updatedOrder.UserUUID,
		PaymentMethod:   updatedOrder.PaymentMethod,
		TransactionUUID: updatedOrder.TransactionUUID,
	}

	err = s.orderPaidProducer.ProduceOrderPaid(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("service: failed to produce OrderPaid event: %w", err)
	}

	return updatedOrder, nil
}
