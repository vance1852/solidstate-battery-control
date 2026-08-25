package worker

import (
	"context"
	"errors"
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
	if IsPermanent(context.DeadlineExceeded) == false {
		t.Fatal("deadline not permanent")
	}
	if IsPermanent(nil) {
		t.Fatal()
	}
}
func TestShouldRetryRejectsDeadline(t *testing.T) {
	p := DefaultRetry()
	if p.ShouldRetry(1, context.DeadlineExceeded) {
		t.Fatal("deadline should not retry")
	}
	if p.ShouldRetry(1, context.Canceled) {
		t.Fatal("cancel should not retry")
	}
	if !p.ShouldRetry(1, errors.New("temporary")) {
		t.Fatal("temporary should retry")
	}
}
