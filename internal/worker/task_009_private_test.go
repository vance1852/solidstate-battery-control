package worker

import (
	"context"
	"testing"
)

func TestCanceledWorkIsNotRetried(t *testing.T) {
	policy := RetryPolicy{MaxAttempts: 4}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	attempts := 0
	err := policy.Run(ctx, func(context.Context) error {
		attempts++
		return context.Canceled
	})
	if attempts != 1 {
		t.Fatalf("canceled work was retried %d times", attempts)
	}
	if err == nil {
		t.Fatal("canceled work unexpectedly succeeded")
	}
}
