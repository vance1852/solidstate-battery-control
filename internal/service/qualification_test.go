package service

import (
	"solidstate-battery-control/internal/domain"
	"testing"
)

func TestQualificationReasons(t *testing.T) {
	w := Workflow{}
	ms := []domain.Measurement{{Kind: "capacity", Value: .1}, {Kind: "resistance", Value: 1}, {Kind: "thermal", Value: 100}, {Kind: "cycle", Value: 1}}
	ok, reasons := w.Qualify(ms)
	if ok || len(reasons) != 4 {
		t.Fatal(ok, reasons)
	}
}
func TestWorkflowReleaseVersion(t *testing.T) {
	w := Workflow{}
	l := domain.CellLot{ID: "lot-1234", State: domain.LotDraft, Version: 1}
	if _, err := w.Release(nil, l, 2); err == nil {
		t.Fatal()
	}
}
