package service

import (
	"context"
	"solidstate-battery-control/internal/domain"
	"testing"
)

func TestPolicyCancellationPaths(t *testing.T) {
	p := Policy{}
	ctx, c := context.WithCancel(context.Background())
	c()
	u := domain.User{Active: true, Role: "admin"}
	l := domain.CellLot{Code: "L"}
	r := domain.QualificationRun{LotID: "L"}
	m := domain.Measurement{Kind: "capacity", Unit: "Ah", Value: 1}
	h := domain.QualityHold{Reason: "x"}
	if p.CheckLot(ctx, u, l) == nil {
		t.Fatal()
	}
	if p.CheckRun(ctx, u, r) == nil {
		t.Fatal()
	}
	if p.CheckMeasurement(ctx, u, m) == nil {
		t.Fatal()
	}
	if p.CheckHold(ctx, u, h) == nil {
		t.Fatal()
	}
}
func TestAuthorizationInactive(t *testing.T) {
	u := domain.User{Active: false, Role: "admin"}
	if Authorize(u, "admin") == nil {
		t.Fatal()
	}
	if RequireOperator(context.Background(), u) == nil {
		t.Fatal()
	}
}
