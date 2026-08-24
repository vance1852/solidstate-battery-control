package domain

import (
	"testing"
	"time"
)

func TestLotTransitions(t *testing.T) {
	cases := []struct {
		from LotState
		to   LotState
		ok   bool
	}{{LotDraft, LotQualified, true}, {LotDraft, LotReleased, false}, {LotQualified, LotReleased, true}, {LotReleased, LotQualified, false}, {LotHold, LotQualified, true}}
	for _, c := range cases {
		if c.from.Can(c.to) != c.ok {
			t.Fatalf("%s to %s", c.from, c.to)
		}
	}
}
func TestRunTransitions(t *testing.T) {
	r := QualificationRun{State: RunPlanned, Version: 1}
	if TransitionRun(&r, RunRunning) != nil {
		t.Fatal("start")
	}
	if TransitionRun(&r, RunSucceeded) != nil {
		t.Fatal("complete")
	}
	if TransitionRun(&r, RunRunning) == nil {
		t.Fatal("terminal transition")
	}
}
func TestHoldTransition(t *testing.T) {
	h := QualityHold{State: HoldOpen}
	if TransitionHold(&h, HoldCleared, "reviewer") != nil {
		t.Fatal("clear")
	}
	if h.ClearedBy == nil || h.ClearedAt == nil {
		t.Fatal("metadata")
	}
	if TransitionHold(&h, HoldOpen, "x") == nil {
		t.Fatal("reopen")
	}
}
func TestRules(t *testing.T) {
	r := DefaultRules()
	if r.Validate() != nil {
		t.Fatal("defaults")
	}
	if !r.AcceptCapacity(1) || r.AcceptCapacity(.1) {
		t.Fatal("capacity")
	}
	if !r.AcceptResistance(.1) || r.AcceptResistance(.3) {
		t.Fatal("resistance")
	}
	if !r.AcceptTemperature(70) || r.AcceptTemperature(90) {
		t.Fatal("temperature")
	}
	if !r.AcceptCycles(100) || r.AcceptCycles(1) {
		t.Fatal("cycles")
	}
}
func TestEvaluate(t *testing.T) {
	now := time.Now()
	ms := []Measurement{{Kind: "capacity", Value: .9, RecordedAt: now}, {Kind: "resistance", Value: .1, RecordedAt: now}, {Kind: "thermal", Value: 70, RecordedAt: now}}
	ok, reasons := DefaultRules().Evaluate(ms)
	if !ok || len(reasons) != 0 {
		t.Fatal(ok, reasons)
	}
	ms[0].Value = .2
	ok, reasons = DefaultRules().Evaluate(ms)
	if ok || len(reasons) == 0 {
		t.Fatal("bad capacity accepted")
	}
}
func TestMeasurementHelpers(t *testing.T) {
	now := time.Now()
	ms := []Measurement{{Kind: "capacity", Value: 1, RecordedAt: now.Add(-time.Hour)}, {Kind: "capacity", Value: .8, RecordedAt: now}, {Kind: "thermal", Value: 70, RecordedAt: now}}
	m, ok := LatestMeasurement(ms, "capacity")
	if !ok || m.Value != .8 {
		t.Fatal("latest")
	}
	if len(MeasurementsWithin(ms, now.Add(-time.Minute), now.Add(time.Minute))) != 2 {
		t.Fatal("window")
	}
	if len(GroupMeasurements(ms)["capacity"]) != 2 {
		t.Fatal("group")
	}
	if Average(ms[:2]) != .9 {
		t.Fatal("average")
	}
	if !HasKind(ms, "thermal") {
		t.Fatal("kind")
	}
}
func TestValidation(t *testing.T) {
	if ValidateLot("", 1) == nil {
		t.Fatal("empty code")
	}
	if ValidateLot("L", 0) == nil {
		t.Fatal("zero capacity")
	}
	if ValidateMeasurement("unknown", "x", 1) == nil {
		t.Fatal("unknown kind")
	}
	if !ValidateRole("admin") || ValidateRole("guest") {
		t.Fatal("role")
	}
	if RequireRole("operator", "reviewer") == nil {
		t.Fatal("role leak")
	}
}
