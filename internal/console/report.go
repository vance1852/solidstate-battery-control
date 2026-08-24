package console

import (
	"solidstate-battery-control/internal/domain"
	"sort"
)

type Report struct {
	Total           int
	Succeeded       int
	Failed          int
	OpenHolds       int
	AverageCapacity float64
}

func BuildReport(lots []domain.CellLot, runs []domain.QualificationRun, holds []domain.QualityHold, ms []domain.Measurement) Report {
	r := Report{Total: len(lots)}
	for _, x := range runs {
		if x.State == domain.RunSucceeded {
			r.Succeeded++
		}
		if x.State == domain.RunFailed {
			r.Failed++
		}
	}
	for _, h := range holds {
		if h.State == domain.HoldOpen {
			r.OpenHolds++
		}
	}
	var values []domain.Measurement
	for _, m := range ms {
		if m.Kind == "capacity" {
			values = append(values, m)
		}
	}
	r.AverageCapacity = domain.Average(values)
	return r
}
func SortRuns(runs []domain.QualificationRun) []domain.QualificationRun {
	out := append([]domain.QualificationRun(nil), runs...)
	sort.Slice(out, func(i, j int) bool { return out[i].ScheduledAt.Before(out[j].ScheduledAt) })
	return out
}
func TerminalRuns(runs []domain.QualificationRun) []domain.QualificationRun {
	out := make([]domain.QualificationRun, 0)
	for _, r := range runs {
		if r.IsComplete() {
			out = append(out, r)
		}
	}
	return out
}
