package service

import (
	"context"
	"solidstate-battery-control/internal/domain"
	"testing"
)

func TestPolicyMatrix(t *testing.T) {
	p := Policy{}
	admin := domain.User{Role: "admin", Active: true}
	operator := domain.User{Role: "operator", Active: true}
	reviewer := domain.User{Role: "reviewer", Active: true}
	lot := domain.CellLot{ID: "lot-1234", Code: "L-1"}
	run := domain.QualificationRun{LotID: lot.ID}
	m := domain.Measurement{Kind: "capacity", Unit: "Ah", Value: 1}
	h := domain.QualityHold{Reason: "drift"}
	if p.CheckLot(context.Background(), operator, lot) != nil {
		t.Fatal("operator lot")
	}
	if p.CheckRun(context.Background(), operator, run) != nil {
		t.Fatal("operator run")
	}
	if p.CheckMeasurement(context.Background(), operator, m) != nil {
		t.Fatal("operator measure")
	}
	if p.CheckHold(context.Background(), reviewer, h) != nil {
		t.Fatal("reviewer hold")
	}
	if p.CheckHold(context.Background(), operator, h) == nil {
		t.Fatal("operator hold")
	}
	if p.CheckLot(context.Background(), admin, lot) != nil {
		t.Fatal("admin lot")
	}
}
