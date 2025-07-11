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
type PaymentMethod int

const (
	UNKNOWN        PaymentMethod = iota // 0
	CARD                                // 1
	SBP                                 // 2
	CREDIT_CARD                         // 3
	INVESTOR_MONEY                      // 4
)

func (pm PaymentMethod) String() string {
	switch pm {
	case CARD:
		return "CARD"
	case SBP:
		return "SBP"
	case CREDIT_CARD:
		return "CREDIT_CARD"
	case INVESTOR_MONEY:
		return "INVESTOR_MONEY"
	default:
		return "UNKNOWN"
	}
}

func PaymentMethodName(value int) string {
	return PaymentMethod(value).String()
}
