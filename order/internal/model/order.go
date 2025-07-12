package model

import "github.com/google/uuid"

type Order struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PartUUIDs       []uuid.UUID
	TotalPrice      float64
	TransactionUUID uuid.UUID
	PaymentMethod   string
	Status          OrderStatus
}
