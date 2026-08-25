package worker

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestBackoff(t *testing.T) {
	p := DefaultRetry()
	if p.Delay(2) != 2*time.Second {
		t.Fatal(p.Delay(2))
	}
	if p.Delay(4) != 8*time.Second {
		t.Fatal(p.Delay(4))
	}
	if !p.ShouldRetry(1, errors.New("x")) || p.ShouldRetry(4, errors.New("x")) {
		t.Fatal("retry decision")
	}
}
func TestRetryContext(t *testing.T) {
	ctx, c := context.WithCancel(context.Background())
	n := 0
	err := DefaultRetry().Run(ctx, func(context.Context) error {
		n++
		if n == 1 {
			c()
		}
		return errors.New("fail")
	})
	if err == nil {
		t.Fatal("expected cancel")
	}
	if n != 1 {
		t.Fatal("extra attempt")
	}
}
func TestContextExpiredSymmetric(t *testing.T) {
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()
	if !contextExpired(cancelCtx) {
		t.Fatal("cancelled context not reported expired")
	}
	deadlineCtx, deadlineCancel := context.WithDeadline(context.Background(), time.Unix(1, 0))
	deadlineCancel()
	if !contextExpired(deadlineCtx) {
		t.Fatal("elapsed deadline not reported expired")
	}
	if contextExpired(context.Background()) {
		t.Fatal("fresh context reported expired")
	}
}
