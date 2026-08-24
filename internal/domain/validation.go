package domain

import (
	"fmt"
	"strings"
)

func ValidateLot(code string, capacity float64) error {
	if strings.TrimSpace(code) == "" {
		return fmt.Errorf("lot code: %w", ErrValidation)
	}
	if capacity <= 0 || capacity > 10000 {
		return fmt.Errorf("capacity: %w", ErrValidation)
	}
	return nil
}
func ValidateMeasurement(kind, unit string, value float64) error {
	allowed := map[string]bool{"capacity": true, "resistance": true, "thermal": true, "cycle": true}
	if !allowed[kind] || strings.TrimSpace(unit) == "" || value < 0 {
		return ErrValidation
	}
	return nil
}
func ValidateRole(role string) bool {
	return role == "admin" || role == "operator" || role == "reviewer"
}
func RequireRole(actual string, allowed ...string) error {
	for _, r := range allowed {
		if actual == r {
			return nil
		}
	}
	return ErrForbidden
}
