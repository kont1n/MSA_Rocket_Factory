package converter

import (
	"github.com/google/uuid"
	"github.com/kont1n/MSA_Rocket_Factory/payment/internal/model"
	paymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
)

func ToModelOrder(req *paymentV1.PayOrderRequest) (model.Order, error) {
	orderUuid, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return model.Order{}, err
	}

	userUuid, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		return model.Order{}, err
	}

	return model.Order{
		OrderUuid:     orderUuid,
		UserUuid:      userUuid,
		PaymentMethod: "CARD", // Используем строковое значение по умолчанию
	}, nil
}
