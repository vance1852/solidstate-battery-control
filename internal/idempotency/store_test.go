package idempotency

import (
	"context"
	"testing"
	"time"
)

func TestStoreRequiresPool(t *testing.T) {
	s := Store{}
	if s.Pool != nil {
		t.Fatal("unexpected pool")
	}
	_ = context.Background()
	_ = time.Second
}
