package model

import (
	"github.com/google/uuid"
)

type Order struct {
	OrderUuid     uuid.UUID
	UserUuid      uuid.UUID
	PaymentMethod string
	TransactionId uuid.UUID
}

// PaymentMethod представляет методы оплаты
type PaymentMethod string

const (
	UNKNOWN        PaymentMethod = "UNKNOWN"
	CARD           PaymentMethod = "CARD"
	SBP            PaymentMethod = "SBP"
	CREDIT_CARD    PaymentMethod = "CREDIT_CARD"
	INVESTOR_MONEY PaymentMethod = "INVESTOR_MONEY"
)

func (pm PaymentMethod) String() string {
	return string(pm)
}

// PaymentMethodFromString создает PaymentMethod из строки
func PaymentMethodFromString(value string) PaymentMethod {
	switch value {
	case "CARD":
		return CARD
	case "SBP":
		return SBP
	case "CREDIT_CARD":
		return CREDIT_CARD
	case "INVESTOR_MONEY":
		return INVESTOR_MONEY
	default:
		return UNKNOWN
	}
}
