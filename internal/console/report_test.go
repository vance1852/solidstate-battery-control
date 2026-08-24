package console

import (
	"solidstate-battery-control/internal/domain"
	"testing"
)

func TestReport(t *testing.T) {
	r := BuildReport([]domain.CellLot{{}, {}}, []domain.QualificationRun{{State: domain.RunSucceeded}, {State: domain.RunFailed}}, []domain.QualityHold{{State: domain.HoldOpen}}, []domain.Measurement{{Kind: "capacity", Value: .9}})
	if r.Total != 2 || r.Succeeded != 1 || r.Failed != 1 || r.OpenHolds != 1 || r.AverageCapacity != .9 {
		t.Fatal(r)
	}
}
