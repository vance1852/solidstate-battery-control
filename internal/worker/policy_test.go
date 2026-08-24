package worker

import (
	"context"
	"testing"
	"time"
)

func TestDelayMonotonic(t *testing.T) {
	p := DefaultRetry()
	last := time.Duration(0)
	for i := 1; i <= p.MaxAttempts; i++ {
		d := p.Delay(i)
		if d <= last {
			t.Fatal(i, d, last)
		}
		last = d
	}
}
func TestPermanent(t *testing.T) {
	if IsPermanent(context.Canceled) == false {
		t.Fatal()
	}
	if IsPermanent(nil) {
		t.Fatal()
	}
}
