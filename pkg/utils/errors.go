package utils

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrForbidden     = errors.New("permission denied")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrBadRequest    = errors.New("bad request")
	ErrInternalError = errors.New("internal server error")
	ErrConflict      = errors.New("resource conflict")
)
