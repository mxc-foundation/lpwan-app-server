package errors

import "errors"

// http errors
var (
	ErrInvalidHeaderName = errors.New("Invalid header name")
)

// influxdb errors
var (
	ErrInvalidPrecision = errors.New("invalid precision value")
)
