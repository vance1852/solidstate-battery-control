package service

import (
	"context"
	"errors"
	"solidstate-battery-control/internal/domain"
	"testing"
	"time"
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
func TestPolicyExpiredDeadlineRejected(t *testing.T) {
	p := Policy{}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	u := domain.User{Active: true, Role: "admin"}
	l := domain.CellLot{Code: "L"}
	r := domain.QualificationRun{LotID: "L"}
	m := domain.Measurement{Kind: "capacity", Unit: "Ah", Value: 1}
	h := domain.QualityHold{Reason: "x"}
	if err := p.CheckContext(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("check context: %v", err)
	}
	if p.CheckLot(ctx, u, l) == nil {
		t.Fatal("expired lot authorized")
	}
	if p.CheckRun(ctx, u, r) == nil {
		t.Fatal("expired run authorized")
	}
	if p.CheckMeasurement(ctx, u, m) == nil {
		t.Fatal("expired measurement authorized")
	}
	if p.CheckHold(ctx, u, h) == nil {
		t.Fatal("expired hold authorized")
	}
	if RequireOperator(ctx, u) == nil {
		t.Fatal("expired operator authorized")
	}
}
func TestPolicyUnexpiredAuthorized(t *testing.T) {
	p := Policy{}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))
	defer cancel()
	u := domain.User{Active: true, Role: "admin"}
	l := domain.CellLot{Code: "L"}
	r := domain.QualificationRun{LotID: "L"}
	m := domain.Measurement{Kind: "capacity", Unit: "Ah", Value: 1}
	h := domain.QualityHold{Reason: "x"}
	if err := p.CheckContext(ctx); err != nil {
		t.Fatalf("check context: %v", err)
	}
	if err := p.CheckLot(ctx, u, l); err != nil {
		t.Fatalf("lot: %v", err)
	}
	if err := p.CheckRun(ctx, u, r); err != nil {
		t.Fatalf("run: %v", err)
	}
	if err := p.CheckMeasurement(ctx, u, m); err != nil {
		t.Fatalf("measurement: %v", err)
	}
	if err := p.CheckHold(ctx, u, h); err != nil {
		t.Fatalf("hold: %v", err)
	}
	if err := RequireOperator(ctx, u); err != nil {
		t.Fatalf("operator: %v", err)
	}
}
