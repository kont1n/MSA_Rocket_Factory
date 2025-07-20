package model

import "errors"

var (
	ErrPartNotFound      = errors.New("part not found")
	ErrPartsListNotFound = errors.New("parts list not found")
	ErrConvertFromRepo   = errors.New("can't parse to model")
	ErrPaid              = errors.New("order status is paid")
	ErrCancelled         = errors.New("order status is cancelled")
	ErrPartsSpecified    = errors.New("parts not specified")
	ErrOrderNotFound      = errors.New("order not found")
)
