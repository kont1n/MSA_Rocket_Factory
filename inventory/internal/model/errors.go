package model

import "errors"

var (
	ErrPartNotFound    = errors.New("part not found")
	ErrConvertFromRepo = errors.New("can't parse to model")
)
