package model

import (
	"github.com/google/uuid"
)

type OrderPaidEvent struct {
	EventUUID       uuid.UUID
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PaymentMethod   string
	TransactionUUID uuid.UUID
}

type ShipAssembledEvent struct {
	EventUUID uuid.UUID
	OrderUUID uuid.UUID
	UserUUID  uuid.UUID
	BuildTime int64
}
