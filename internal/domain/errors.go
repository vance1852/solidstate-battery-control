package domain

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrConflict     = errors.New("concurrent update conflict")
	ErrForbidden    = errors.New("operation forbidden")
	ErrUnauthorized = errors.New("authentication required")
	ErrInvalidState = errors.New("invalid state transition")
	ErrValidation   = errors.New("validation failed")
)
