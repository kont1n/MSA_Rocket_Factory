package model

import "errors"

var (
	ErrPartsListNotFound = errors.New("parts list not found")
	ErrConvertFromRepo   = errors.New("can't parse to model")
	ErrConvertFromClient = errors.New("can't parse to model")
	ErrInventoryClient   = errors.New("inventory client error")
	ErrPaymentClient     = errors.New("payment client error")
	ErrPaid              = errors.New("order status is paid")
	ErrCancelled         = errors.New("order status is cancelled")
	ErrPartsSpecified    = errors.New("parts not specified")
	ErrOrderNotFound     = errors.New("order not found")
)
