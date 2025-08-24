package kafka

import "github.com/kont1n/MSA_Rocket_Factory/order/internal/model"

type ShipAssembledDecoder interface {
	Decode(data []byte) (*model.ShipAssembledEvent, error)
}
