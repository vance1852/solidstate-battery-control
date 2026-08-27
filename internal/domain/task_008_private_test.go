package domain

import "testing"

func TestClearedHoldCannotBeClearedAgain(t *testing.T) {
	user := "reviewer-1"
	hold := QualityHold{ID: "hold-008", State: HoldCleared}
	if err := TransitionHold(&hold, HoldCleared, user); err == nil {
		t.Fatal("a cleared hold was accepted for a second clear transition")
	}
}
