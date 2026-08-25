package domain

import "testing"

func TestLotTransitionFailurePreservesVersion(t *testing.T) {
	lot := CellLot{ID: "lot-001", State: LotDraft, Version: 7}
	err := TransitionLot(&lot, LotReleased)
	if err == nil {
		t.Fatal("draft lot was released without qualification")
	}
	if lot.State != LotDraft {
		t.Fatalf("failed transition changed state to %s", lot.State)
	}
	if lot.Version != 7 {
		t.Fatalf("failed transition consumed version %d", lot.Version)
	}
	qualified := CellLot{ID: "lot-002", State: LotDraft, Version: 3}
	if err := TransitionLot(&qualified, LotQualified); err != nil {
		t.Fatal(err)
	}
	if qualified.State != LotQualified || qualified.Version != 4 {
		t.Fatalf("valid transition not committed: %+v", qualified)
	}
}
