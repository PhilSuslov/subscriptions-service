package subscription

import "errors"

var (
	ErrInvalidServiceName = errors.New("service name must not be empty")
	ErrInvalidPrice       = errors.New("price must be positive")
	ErrInvalidPeriod      = errors.New("end month must be greater than or equal to start month")
	ErrNotFound           = errors.New("subscription not found")
	ErrInvalidUserID      = errors.New("user id must be valid uuid")
)
