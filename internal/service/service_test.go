package service

import (
	"context"
	"solidstate-battery-control/internal/domain"
	"testing"
)

func TestAuthorizeRoles(t *testing.T) {
	u := domain.User{Active: true, Role: "operator"}
	if Authorize(u, "reviewer") == nil {
		t.Fatal("operator must not review")
	}
	if Authorize(u, "operator") != nil {
		t.Fatal("operator denied")
	}
}
func TestContextCancellation(t *testing.T) {
	ctx, c := context.WithCancel(context.Background())
	c()
	if RequireOperator(ctx, domain.User{Active: true, Role: "operator"}) == nil {
		t.Fatal("cancel ignored")
	}
}
