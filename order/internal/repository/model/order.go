package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderUUID       string
	UserUUID        string
	PartUUIDs       []string
	TotalPrice      float32
	TransactionUUID string
	PaymentMethod   string
	Status          string
}

type OrderPostgres struct {
	OrderUUID       uuid.UUID   `db:"order_uuid"`
	UserUUID        uuid.UUID   `db:"user_uuid"`
	PartUUIDs       []uuid.UUID `db:"part_uuid"`
	TotalPrice      float32     `db:"total_price"`
	TransactionUUID uuid.UUID   `db:"transaction_uuid"`
	PaymentMethod   string      `db:"payment_method"`
	Status          string      `db:"status"`
	CreatedAt       time.Time   `db:"created_at"`
	UpdatedAt       time.Time   `db:"updated_at"`
}
