package model

type OrderPaid struct {
	EventUUID       string
	OrderUUID       string
	UserUUID        string
	PaymentMethod   string
	TransactionUUID string
}
