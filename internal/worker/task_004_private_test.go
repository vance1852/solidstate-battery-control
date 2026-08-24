package worker

import (
	"context"
	"testing"
	"time"
)

func TestDeadlineStopsRetryBeforeDelivery(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Unix(0, 0))
	defer cancel()
	if !contextExpired(ctx) {
		t.Fatal("deadline was not recognized")
	}
	deliveries := 0
	err := DefaultRetry().Run(ctx, func(context.Context) error {
		deliveries++
		return nil
	})
	if err == nil {
		t.Fatal("expired worker accepted a delivery")
	}
	if deliveries != 0 {
		t.Fatalf("expired worker delivered %d times", deliveries)
	}
}
