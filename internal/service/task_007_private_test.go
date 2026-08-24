package service

import (
	"solidstate-battery-control/internal/domain"
	"testing"
)

func TestUnknownRoleCannotUseOperatorCapability(t *testing.T) {
	user := domain.User{ID: "u-unknown", Role: "quality_auditor", Active: true}
	if err := Authorize(user, "operator", "admin"); err == nil {
		t.Fatal("unknown role was granted an operator capability")
	}
}
