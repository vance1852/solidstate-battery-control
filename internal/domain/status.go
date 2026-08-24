package domain

type LotState string

const (
	LotDraft     LotState = "draft"
	LotQualified LotState = "qualified"
	LotReleased  LotState = "released"
	LotHold      LotState = "hold"
)

type RunState string

const (
	RunPlanned   RunState = "planned"
	RunRunning   RunState = "running"
	RunPaused    RunState = "paused"
	RunSucceeded RunState = "succeeded"
	RunFailed    RunState = "failed"
)

type HoldState string

const (
	HoldOpen    HoldState = "open"
	HoldCleared HoldState = "cleared"
)

func (s LotState) IsTerminal() bool {
	return s == LotReleased
}

func (s LotState) Can(next LotState) bool {
	switch s {
	case LotDraft:
		return next == LotQualified || next == LotHold
	case LotQualified:
		return next == LotReleased || next == LotHold
	case LotHold:
		return next == LotQualified || next == LotReleased
	case LotReleased:
		return next == LotHold
	}
	return false
}
func (s RunState) Can(next RunState) bool {
	switch s {
	case RunPlanned:
		return next == RunRunning || next == RunFailed
	case RunRunning:
		return next == RunPaused || next == RunSucceeded || next == RunFailed
	case RunPaused:
		return next == RunRunning || next == RunFailed
	}
	return false
}
func (s HoldState) Can(next HoldState) bool { return s == HoldOpen && next == HoldCleared }
