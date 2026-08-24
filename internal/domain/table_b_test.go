package domain

import "testing"

func TestTableB01(t *testing.T) {
	if !DefaultRules().AcceptCapacity(1) {
		t.Fatal()
	}
}
func TestTableB02(t *testing.T) {
	if DefaultRules().AcceptCapacity(0) {
		t.Fatal()
	}
}
func TestTableB03(t *testing.T) {
	if !DefaultRules().AcceptResistance(.1) {
		t.Fatal()
	}
}
func TestTableB04(t *testing.T) {
	if DefaultRules().AcceptResistance(1) {
		t.Fatal()
	}
}
func TestTableB05(t *testing.T) {
	if !DefaultRules().AcceptTemperature(50) {
		t.Fatal()
	}
}
func TestTableB06(t *testing.T) {
	if DefaultRules().AcceptTemperature(100) {
		t.Fatal()
	}
}
func TestTableB07(t *testing.T) {
	if !DefaultRules().AcceptCycles(100) {
		t.Fatal()
	}
}
func TestTableB08(t *testing.T) {
	if DefaultRules().AcceptCycles(10) {
		t.Fatal()
	}
}
func TestTableB09(t *testing.T) {
	if NormalizeRole("") != "operator" {
		t.Fatal()
	}
}
func TestTableB10(t *testing.T) {
	if NormalizeState("") != "draft" {
		t.Fatal()
	}
}
func TestTableB11(t *testing.T) {
	if !(CellLot{State: LotQualified}).CanSchedule() {
		t.Fatal()
	}
}
func TestTableB12(t *testing.T) {
	if (CellLot{State: LotDraft}).CanSchedule() {
		t.Fatal()
	}
}
func TestTableB13(t *testing.T) {
	if !(QualificationRun{State: RunRunning}).IsActive() {
		t.Fatal()
	}
}
func TestTableB14(t *testing.T) {
	if (QualificationRun{State: RunSucceeded}).IsActive() {
		t.Fatal()
	}
}
func TestTableB15(t *testing.T) {
	if !(QualificationRun{State: RunSucceeded}).IsComplete() {
		t.Fatal()
	}
}
func TestTableB16(t *testing.T) {
	if (QualificationRun{State: RunRunning}).IsComplete() {
		t.Fatal()
	}
}
func TestTableB17(t *testing.T) {
	if !(QualityHold{State: HoldOpen}).IsOpen() {
		t.Fatal()
	}
}
func TestTableB18(t *testing.T) {
	if (QualityHold{State: HoldCleared}).IsOpen() {
		t.Fatal()
	}
}
func TestTableB19(t *testing.T) {
	if JoinErrors(nil, nil) != nil {
		t.Fatal()
	}
}
func TestTableB20(t *testing.T) {
	if JoinErrors(ErrConflict, nil) != ErrConflict {
		t.Fatal()
	}
}
