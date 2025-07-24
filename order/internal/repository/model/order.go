package model

type Order struct {
	OrderUUID       string
	UserUUID        string
	PartUUIDs       []string
	TotalPrice      float32
	TransactionUUID string
	PaymentMethod   string
	Status          string
}
