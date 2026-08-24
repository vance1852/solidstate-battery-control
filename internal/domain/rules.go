package domain

import (
	"fmt"
	"strings"
)

type RuleSet struct {
	MaxTemperature float64
	MinCapacity    float64
	MaxResistance  float64
	RequiredCycles int
}

func DefaultRules() RuleSet {
	return RuleSet{MaxTemperature: 80, MinCapacity: 0.8, MaxResistance: 0.2, RequiredCycles: 100}
}
func (r RuleSet) Validate() error {
	if r.MaxTemperature <= 0 {
		return fmt.Errorf("temperature must be positive")
	}
	if r.MinCapacity <= 0 || r.MinCapacity > 1 {
		return fmt.Errorf("capacity threshold invalid")
	}
	if r.MaxResistance <= 0 {
		return fmt.Errorf("resistance threshold invalid")
	}
	if r.RequiredCycles < 1 {
		return fmt.Errorf("cycles required")
	}
	return nil
}
func (r RuleSet) AcceptCapacity(v float64) bool    { return v >= r.MinCapacity }
func (r RuleSet) AcceptResistance(v float64) bool  { return v <= r.MaxResistance }
func (r RuleSet) AcceptTemperature(v float64) bool { return v <= r.MaxTemperature }
func (r RuleSet) AcceptCycles(v int) bool          { return v >= r.RequiredCycles }
func (r RuleSet) Evaluate(m []Measurement) (bool, []string) {
	ok := true
	reasons := make([]string, 0)
	for _, x := range m {
		if !MeasurementKindKnown(x.Kind) {
			ok = false
			reasons = append(reasons, "unknown measurement kind")
			continue
		}
		switch x.Kind {
		case "capacity":
			if !r.AcceptCapacity(x.Value) {
				ok = false
				reasons = append(reasons, "capacity below threshold")
			}
		case "resistance":
			if !r.AcceptResistance(x.Value) {
				ok = false
				reasons = append(reasons, "resistance above threshold")
			}
		case "thermal":
			if !r.AcceptTemperature(x.Value) {
				ok = false
				reasons = append(reasons, "temperature above threshold")
			}
		case "cycle":
			if !r.AcceptCycles(int(x.Value)) {
				ok = false
				reasons = append(reasons, "cycle count incomplete")
			}
		}
	}
	return ok, reasons
}
func (l CellLot) IsTerminal() bool  { return l.State == LotReleased }
func (l CellLot) CanSchedule() bool { return l.State == LotQualified || l.State == LotReleased }
func (r QualificationRun) IsActive() bool {
	return r.State == RunPlanned || r.State == RunRunning || r.State == RunPaused
}
func (r QualificationRun) IsComplete() bool { return r.State == RunSucceeded || r.State == RunFailed }
func (h QualityHold) IsOpen() bool          { return h.State == HoldOpen }
func NormalizeRole(role string) string {
	role = strings.ToLower(strings.TrimSpace(role))
	if !ValidateRole(role) {
		return "operator"
	}
	return role
}
func NormalizeState(state string) string {
	if state == "" {
		return "draft"
	}
	return state
}
func EnsureVersion(v int64) error {
	if v < 1 {
		return ErrConflict
	}
	return nil
}
func EnsurePositive(v float64) error {
	if v <= 0 {
		return ErrValidation
	}
	return nil
}
func EnsureID(id string) error {
	if len(id) < 8 {
		return ErrValidation
	}
	return nil
}
func JoinErrors(a, b error) error {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	return fmt.Errorf("%w; %v", a, b)
}
