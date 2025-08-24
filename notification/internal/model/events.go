package model

import (
	"time"

	"github.com/google/uuid"
)

// OrderPaidEvent - событие оплаты заказа
type OrderPaidEvent struct {
	EventUUID       uuid.UUID
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PaymentMethod   string
	TransactionUUID uuid.UUID
}

// ShipAssembledEvent - событие сборки корабля
type ShipAssembledEvent struct {
	EventUUID uuid.UUID
	OrderUUID uuid.UUID
	UserUUID  uuid.UUID
	BuildTime int64
}

// SightingInfo - информация о наблюдении UFO
type SightingInfo struct {
	Location        string
	Description     string
	ObservedAt      *time.Time
	Color           *string
	Sound           *bool
	DurationSeconds *int32
}
