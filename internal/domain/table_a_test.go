package domain

import "testing"

func TestTableA01(t *testing.T) {
	if !LotDraft.Can(LotQualified) {
		t.Fatal()
	}
}
func TestTableA02(t *testing.T) {
	if LotDraft.Can(LotReleased) {
		t.Fatal()
	}
}
func TestTableA03(t *testing.T) {
	if !LotQualified.Can(LotReleased) {
		t.Fatal()
	}
}
func TestTableA04(t *testing.T) {
	if LotQualified.Can(LotDraft) {
		t.Fatal()
	}
}
func TestTableA05(t *testing.T) {
	if !LotReleased.Can(LotHold) {
		t.Fatal()
	}
}
func TestTableA06(t *testing.T) {
	if LotReleased.Can(LotDraft) {
		t.Fatal()
	}
}
func TestTableA07(t *testing.T) {
	if !LotHold.Can(LotQualified) {
		t.Fatal()
	}
}
func TestTableA08(t *testing.T) {
	if !LotHold.Can(LotReleased) {
		t.Fatal()
	}
}
func TestTableA09(t *testing.T) {
	if LotHold.Can(LotDraft) {
		t.Fatal()
	}
}
func TestTableA10(t *testing.T) {
	if !RunPlanned.Can(RunRunning) {
		t.Fatal()
	}
}
func TestTableA11(t *testing.T) {
	if !RunPlanned.Can(RunFailed) {
		t.Fatal()
	}
}
func TestTableA12(t *testing.T) {
	if RunPlanned.Can(RunSucceeded) {
		t.Fatal()
	}
}
func TestTableA13(t *testing.T) {
	if !RunRunning.Can(RunPaused) {
		t.Fatal()
	}
}
func TestTableA14(t *testing.T) {
	if !RunRunning.Can(RunSucceeded) {
		t.Fatal()
	}
}
func TestTableA15(t *testing.T) {
	if !RunRunning.Can(RunFailed) {
		t.Fatal()
	}
}
func TestTableA16(t *testing.T) {
	if !RunPaused.Can(RunRunning) {
		t.Fatal()
	}
}
func TestTableA17(t *testing.T) {
	if !RunPaused.Can(RunFailed) {
		t.Fatal()
	}
}
func TestTableA18(t *testing.T) {
	if RunPaused.Can(RunSucceeded) {
		t.Fatal()
	}
}
func TestTableA19(t *testing.T) {
	if RunSucceeded.Can(RunRunning) {
		t.Fatal()
	}
}
func TestTableA20(t *testing.T) {
	if RunFailed.Can(RunRunning) {
		t.Fatal()
	}
}
func TestTableA21(t *testing.T) {
	if !HoldOpen.Can(HoldCleared) {
		t.Fatal()
	}
}
func TestTableA22(t *testing.T) {
	if HoldCleared.Can(HoldOpen) {
		t.Fatal()
	}
}
func TestTableA23(t *testing.T) {
	if !ValidateRole("admin") {
		t.Fatal()
	}
}
func TestTableA24(t *testing.T) {
	if ValidateRole("guest") {
		t.Fatal()
	}
}
func TestTableA25(t *testing.T) {
	if EnsurePositive(1) != nil {
		t.Fatal()
	}
}
func TestTableA26(t *testing.T) {
	if EnsurePositive(0) == nil {
		t.Fatal()
	}
}
func TestTableA27(t *testing.T) {
	if EnsureVersion(1) != nil {
		t.Fatal()
	}
}
func TestTableA28(t *testing.T) {
	if EnsureVersion(0) == nil {
		t.Fatal()
	}
}
func TestTableA29(t *testing.T) {
	if EnsureID("12345678") != nil {
		t.Fatal()
	}
}
func TestTableA30(t *testing.T) {
	if EnsureID("x") == nil {
		t.Fatal()
	}
}
