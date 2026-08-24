package domain

import (
	"errors"
	"testing"
)

func TestConflictSentinelSurvivesServiceWrapping(t *testing.T) {
	err := WrapOperation("change lot state", ErrConflict)
	if err == nil {
		t.Fatal("missing wrapped error")
	}
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("conflict identity was lost: %v", err)
	}
	if errors.Is(WrapOperation("validation", ErrValidation), ErrConflict) {
		t.Fatal("unrelated sentinel matched")
	}
}
