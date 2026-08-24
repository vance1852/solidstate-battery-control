package domain

import "testing"

func TestUnknownMeasurementKindCannotQualifyRun(t *testing.T) {
	measurements := []Measurement{
		{Kind: "capacity", Value: 0.9},
		{Kind: "resistance", Value: 0.1},
		{Kind: "mystery", Value: 42},
	}
	ok, reasons := DefaultRules().Evaluate(measurements)
	if ok {
		t.Fatal("run with unknown measurement was accepted")
	}
	if len(reasons) == 0 {
		t.Fatal("quality rejection did not explain the unknown measurement")
	}
	if MeasurementKindKnown("mystery") {
		t.Fatal("unknown measurement was classified as known")
	}
}
