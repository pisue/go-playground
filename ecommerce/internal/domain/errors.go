package domain

import "errors"

var (
	ErrBadRequest      = errors.New("bad_request")
	ErrNotFound        = errors.New("not_found")
	ErrInternalFailure = errors.New("internal_failure")
)

type AppError struct {
	ServiceError error
	Detail       error
	Message      string
}

func (e *AppError) Error() string {
	return e.Message
}
