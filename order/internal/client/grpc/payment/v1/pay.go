package v1

import (
	"context"
	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/client/converter"
	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	generaredPaymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
)

func (p paymentClient) CreatePayment(ctx context.Context, order *model.Order) (*model.Order, error) {
	// Оплачиваем заказ с помощью gRPC клиента
	response, err := p.generatedClient.PayOrder(ctx, &generaredPaymentV1.PayOrderRequest{
		OrderUuid:     order.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: converter.PaymentToProto(order.PaymentMethod),
	})
	if err != nil {
		return nil, err
	}

	transaction, err := uuid.Parse(response.GetTransactionUuid())
	if err != nil {
		return nil, err
	}

	// Обновляем заказ
	order.TransactionUUID = transaction
	order.Status = model.StatusPaid

	return order, nil
}
