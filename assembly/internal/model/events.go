package model

import "time"

type AssemblyRecordedEvent struct {
	UUID        string
	ObservedAt  *time.Time
	Location    string
	Description string
}

type OrderPaid struct {
	EventUUID       string
	OrderUUID       string
	UserUUID        string
	PaymentMethod   string
	TransactionUUID string
}

type ShipAssembled struct {
	EventUUID string
	OrderUUID string
	UserUUID  string
	BuildTime int64
}
