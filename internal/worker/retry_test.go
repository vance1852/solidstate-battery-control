package worker

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryEventuallySucceeds(t *testing.T) {
	n := 0
	err := DefaultRetry().Run(context.Background(), func(context.Context) error {
		n++
		if n < 3 {
			return errors.New("temporary")
		}
		return nil
	})
	if err != nil || n != 3 {
		t.Fatalf("err=%v attempts=%d", err, n)
	}
}
func TestRetryCancellation(t *testing.T) {
	ctx, c := context.WithCancel(context.Background())
	c()
	if DefaultRetry().Run(ctx, func(context.Context) error { return errors.New("x") }) == nil {
		t.Fatal("cancellation ignored")
	}
}
func TestRetryDeadlineStopsDelivery(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Unix(1, 0))
	defer cancel()
	n := 0
	err := DefaultRetry().Run(ctx, func(context.Context) error {
		n++
		return errors.New("x")
	})
	if err == nil {
		t.Fatal("deadline ignored")
	}
	if n != 0 {
		t.Fatalf("expired task received delivery callback: %d", n)
	}
}
func TestRetryDeliversWhileNotExpired(t *testing.T) {
	n := 0
	err := DefaultRetry().Run(context.Background(), func(context.Context) error {
		n++
		if n < 2 {
			return errors.New("temporary")
		}
		return nil
	})
	if err != nil || n != 2 {
		t.Fatalf("err=%v attempts=%d", err, n)
	}
}
