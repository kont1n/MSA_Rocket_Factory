package model

import "errors"

var (
	ErrConvertFromKafkaEvent = errors.New("can't parse to model")
	ErrMarshalToKafkaEvent   = errors.New("failed to marshal ShipAssembled")
	ErrSendToKafka           = errors.New("failed to publish ShipAssembled")
)
