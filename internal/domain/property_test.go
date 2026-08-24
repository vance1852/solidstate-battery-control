package domain

import "testing"

func TestStateTableOne(t *testing.T) {
	if !LotDraft.Can(LotQualified) {
		t.Fatal()
	}
	if LotDraft.Can(LotReleased) {
		t.Fatal()
	}
	if !LotQualified.Can(LotReleased) {
		t.Fatal()
	}
	if LotQualified.Can(LotDraft) {
		t.Fatal()
	}
	if !LotReleased.Can(LotHold) {
		t.Fatal()
	}
	if LotReleased.Can(LotDraft) {
		t.Fatal()
	}
}
func TestStateTableTwo(t *testing.T) {
	if !LotHold.Can(LotQualified) {
		t.Fatal()
	}
	if !LotHold.Can(LotReleased) {
		t.Fatal()
	}
	if LotHold.Can(LotDraft) {
		t.Fatal()
	}
	if !RunPlanned.Can(RunRunning) {
		t.Fatal()
	}
	if !RunPlanned.Can(RunFailed) {
		t.Fatal()
	}
	if RunPlanned.Can(RunSucceeded) {
		t.Fatal()
	}
}
func TestStateTableThree(t *testing.T) {
	if !RunRunning.Can(RunPaused) {
		t.Fatal()
	}
	if !RunRunning.Can(RunSucceeded) {
		t.Fatal()
	}
	if !RunRunning.Can(RunFailed) {
		t.Fatal()
	}
	if !RunPaused.Can(RunRunning) {
		t.Fatal()
	}
	if !RunPaused.Can(RunFailed) {
		t.Fatal()
	}
	if RunPaused.Can(RunSucceeded) {
		t.Fatal()
	}
}
func TestStateTableFour(t *testing.T) {
	if RunSucceeded.Can(RunRunning) {
		t.Fatal()
	}
	if RunFailed.Can(RunRunning) {
		t.Fatal()
	}
	if HoldCleared.Can(HoldOpen) {
		t.Fatal()
	}
	if !HoldOpen.Can(HoldCleared) {
		t.Fatal()
	}
}
func TestRuleBoundariesOne(t *testing.T) {
	r := DefaultRules()
	if r.AcceptCapacity(r.MinCapacity - 0.01) {
		t.Fatal()
	}
	if !r.AcceptCapacity(r.MinCapacity) {
		t.Fatal()
	}
	if r.AcceptResistance(r.MaxResistance + 0.01) {
		t.Fatal()
	}
	if !r.AcceptResistance(r.MaxResistance) {
		t.Fatal()
	}
}
func TestRuleBoundariesTwo(t *testing.T) {
	r := DefaultRules()
	if r.AcceptTemperature(r.MaxTemperature + 1) {
		t.Fatal()
	}
	if !r.AcceptTemperature(r.MaxTemperature) {
		t.Fatal()
	}
	if r.AcceptCycles(r.RequiredCycles - 1) {
		t.Fatal()
	}
	if !r.AcceptCycles(r.RequiredCycles) {
		t.Fatal()
	}
}
func TestValidationCasesOne(t *testing.T) {
	for _, code := range []string{"", " ", "\t"} {
		if ValidateLot(code, 1) == nil {
			t.Fatal(code)
		}
	}
	for _, cap := range []float64{0, -1, 10001} {
		if ValidateLot("L", cap) == nil {
			t.Fatal(cap)
		}
	}
}
func TestValidationCasesTwo(t *testing.T) {
	for _, kind := range []string{"x", "", "voltage"} {
		if ValidateMeasurement(kind, "u", 1) == nil {
			t.Fatal(kind)
		}
	}
	if ValidateMeasurement("capacity", "", 1) == nil {
		t.Fatal("unit")
	}
	if ValidateMeasurement("capacity", "Ah", -1) == nil {
		t.Fatal("negative")
	}
}
func TestRoleCases(t *testing.T) {
	for _, role := range []string{"admin", "operator", "reviewer"} {
		if !ValidateRole(role) {
			t.Fatal(role)
		}
	}
	for _, role := range []string{"guest", "", "root"} {
		if ValidateRole(role) {
			t.Fatal(role)
		}
	}
}
func TestErrorJoin(t *testing.T) {
	if JoinErrors(nil, nil) != nil {
		t.Fatal()
	}
	if JoinErrors(ErrConflict, nil) != ErrConflict {
		t.Fatal()
	}
	if JoinErrors(nil, ErrValidation) != ErrValidation {
		t.Fatal()
	}
}
