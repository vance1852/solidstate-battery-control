package domain

import (
	"testing"
	"time"
)

func TestMeasurementKinds(t *testing.T) {
	ms := []Measurement{{Kind: "capacity", Value: .9}, {Kind: "resistance", Value: .1}, {Kind: "thermal", Value: 70}, {Kind: "cycle", Value: 100}}
	for _, kind := range []string{"capacity", "resistance", "thermal", "cycle"} {
		if !HasKind(ms, kind) {
			t.Fatal(kind)
		}
	}
	if HasKind(ms, "unknown") {
		t.Fatal("unknown")
	}
}
func TestAverageEmpty(t *testing.T) {
	if Average(nil) != 0 {
		t.Fatal()
	}
	if Average([]Measurement{}) != 0 {
		t.Fatal()
	}
}
func TestWindowBoundaries(t *testing.T) {
	t0 := time.Now()
	ms := []Measurement{{RecordedAt: t0}, {RecordedAt: t0.Add(time.Second)}, {RecordedAt: t0.Add(2 * time.Second)}}
	if len(MeasurementsWithin(ms, t0, t0.Add(2*time.Second))) != 2 {
		t.Fatal()
	}
	if len(MeasurementsWithin(ms, t0.Add(2*time.Second), t0.Add(3*time.Second))) != 1 {
		t.Fatal()
	}
}
func TestLatestMissing(t *testing.T) {
	if _, ok := LatestMeasurement(nil, "capacity"); ok {
		t.Fatal()
	}
	if _, ok := LatestMeasurement([]Measurement{{Kind: "thermal"}}, "capacity"); ok {
		t.Fatal()
	}
}
func TestGrouping(t *testing.T) {
	g := GroupMeasurements([]Measurement{{Kind: "capacity"}, {Kind: "capacity"}, {Kind: "thermal"}})
	if len(g) != 2 || len(g["capacity"]) != 2 || len(g["thermal"]) != 1 {
		t.Fatal(g)
	}
}
