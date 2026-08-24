package clock

import (
	"testing"
	"time"
)

func TestFixedClock(t *testing.T) {
	want := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if (Fixed{T: want}).Now() != want {
		t.Fatal("clock")
	}
	if (Real{}).Now().IsZero() {
		t.Fatal("real clock")
	}
}
