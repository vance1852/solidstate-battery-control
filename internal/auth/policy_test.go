package auth

import (
	"testing"
	"time"
)

func TestPolicy(t *testing.T) {
	if NormalizeEmail(" A@B ") != "a@b" {
		t.Fatal("email")
	}
	if NormalizeRole("guest") != "operator" {
		t.Fatal("role")
	}
	now := time.Now()
	if !SessionValid(now.Add(time.Hour), nil, now) {
		t.Fatal("valid")
	}
	if SessionValid(now.Add(-time.Hour), nil, now) {
		t.Fatal("expired")
	}
	if !Can("operator", "lot:create") || Can("operator", "user:manage") {
		t.Fatal("actions")
	}
	if len(Actions("admin")) != 7 {
		t.Fatal("admin actions")
	}
}
