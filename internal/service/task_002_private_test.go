package service

import (
	"context"
	"testing"
	"time"
	"solidstate-battery-control/internal/domain"
)

func TestDeadlineCancellationReachesAuthorization(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Unix(0, 0))
	defer cancel()
	user := domain.User{Active: true, Role: "operator"}
	if err := RequireOperator(ctx, user); err == nil {
		t.Fatal("expired request was authorized")
	}
	if err := (Policy{}).CheckContext(ctx); err == nil {
		t.Fatal("expired request reached policy without cancellation")
	}
}
