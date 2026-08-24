package domain

import (
	"testing"
	"time"
)

func TestFilterAndPaging(t *testing.T) {
	lots := []CellLot{{Code: "A-1", State: LotDraft, UpdatedAt: time.Now()}, {Code: "B-2", State: LotQualified, UpdatedAt: time.Now().Add(time.Hour)}, {Code: "A-3", State: LotHold, UpdatedAt: time.Now().Add(2 * time.Hour)}}
	f := LotFilter{Query: "a", States: []LotState{LotDraft}}
	if !f.Match(lots[0]) || f.Match(lots[1]) {
		t.Fatal("filter")
	}
	if len(PageLots(lots, 2, 1)) != 2 {
		t.Fatal("page")
	}
	if len(PageLots(lots, 2, 10)) != 0 {
		t.Fatal("bounds")
	}
	if len(SortLots(lots, true)) != 3 {
		t.Fatal("sort")
	}
	if len(ParseStates([]string{"draft", "bad", "hold"})) != 2 {
		t.Fatal("states")
	}
}
func TestLotProperties(t *testing.T) {
	l := CellLot{State: LotQualified, ID: "12345678"}
	if !l.CanSchedule() || l.IsTerminal() {
		t.Fatal("properties")
	}
	if EnsureID(l.ID) != nil || EnsureID("x") == nil {
		t.Fatal("id")
	}
	if EnsureVersion(0) == nil || EnsureVersion(1) != nil {
		t.Fatal("version")
	}
	if EnsurePositive(0) == nil || EnsurePositive(1) != nil {
		t.Fatal("positive")
	}
}
