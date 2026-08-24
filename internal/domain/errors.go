package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound     = errors.New("resource not found")
	ErrConflict     = errors.New("concurrent update conflict")
	ErrForbidden    = errors.New("operation forbidden")
	ErrUnauthorized = errors.New("authentication required")
	ErrInvalidState = errors.New("invalid state transition")
	ErrValidation   = errors.New("validation failed")
)

func WrapOperation(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %v", op, err)
}
