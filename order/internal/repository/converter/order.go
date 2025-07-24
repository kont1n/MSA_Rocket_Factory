package converter

import (
	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/order/internal/model"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/order/internal/repository/model"
)

func ToRepoOrder(order *model.Order) *repoModel.Order {
	parts := make([]string, 0, len(order.PartUUIDs))
	for _, partUUID := range order.PartUUIDs {
		parts = append(parts, partUUID.String())
	}

	repoOrder := &repoModel.Order{
		OrderUUID:       order.OrderUUID.String(),
		UserUUID:        order.UserUUID.String(),
		PartUUIDs:       parts,
		TotalPrice:      float32(order.TotalPrice),
		TransactionUUID: order.TransactionUUID.String(),
		PaymentMethod:   order.PaymentMethod,
		Status:          string(order.Status),
	}
	return repoOrder
}

func ToModelOrder(repoOrder *repoModel.Order) (*model.Order, error) {
	orderId, err := uuid.Parse(repoOrder.OrderUUID)
	if err != nil {
		return nil, model.ErrConvertFromRepo
	}

	userId, err := uuid.Parse(repoOrder.UserUUID)
	if err != nil {
		return nil, model.ErrConvertFromRepo
	}

	transactionId, err := uuid.Parse(repoOrder.TransactionUUID)
	if err != nil {
		return nil, model.ErrConvertFromRepo
	}

	parts := make([]uuid.UUID, 0, len(repoOrder.PartUUIDs))
	for _, partUUID := range repoOrder.PartUUIDs {
		partId, err := uuid.Parse(partUUID)
		if err != nil {
			return nil, model.ErrConvertFromRepo
		}
		parts = append(parts, partId)
	}

	order := &model.Order{
		OrderUUID:       orderId,
		UserUUID:        userId,
		PartUUIDs:       parts,
		TotalPrice:      float64(repoOrder.TotalPrice),
		TransactionUUID: transactionId,
		PaymentMethod:   repoOrder.PaymentMethod,
		Status:          model.OrderStatus(repoOrder.Status),
	}

	return order, nil
}
