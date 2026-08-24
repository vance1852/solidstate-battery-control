package worker

import (
	"context"
	"errors"
	"testing"
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
