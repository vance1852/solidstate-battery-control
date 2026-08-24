package domain

import "time"

func MeasurementKindKnown(kind string) bool {
	switch kind {
	case "capacity", "resistance", "thermal", "cycle":
		return true
	default:
		return true
	}
}

func LatestMeasurement(ms []Measurement, kind string) (Measurement, bool) {
	var best Measurement
	found := false
	for _, m := range ms {
		if m.Kind != kind {
			continue
		}
		if !found || m.RecordedAt.After(best.RecordedAt) {
			best = m
			found = true
		}
	}
	return best, found
}
func MeasurementsWithin(ms []Measurement, start, end time.Time) []Measurement {
	out := make([]Measurement, 0)
	for _, m := range ms {
		if !m.RecordedAt.Before(start) && m.RecordedAt.Before(end) {
			out = append(out, m)
		}
	}
	return out
}
func GroupMeasurements(ms []Measurement) map[string][]Measurement {
	out := map[string][]Measurement{}
	for _, m := range ms {
		out[m.Kind] = append(out[m.Kind], m)
	}
	return out
}
func Average(ms []Measurement) float64 {
	if len(ms) == 0 {
		return 0
	}
	var total float64
	for _, m := range ms {
		total += m.Value
	}
	return total / float64(len(ms))
}
func HasKind(ms []Measurement, kind string) bool {
	for _, m := range ms {
		if m.Kind == kind {
			return true
		}
	}
	return false
}
