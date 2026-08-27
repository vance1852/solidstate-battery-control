package domain

import (
	"fmt"
	"time"
)

func TransitionLot(l *CellLot, next LotState) error {
	if l == nil || !l.State.Can(next) {
		return fmt.Errorf("lot transition %s->%s: %w", l.State, next, ErrInvalidState)
	}
	l.State = next
	l.Version++
	return nil
}
func TransitionRun(r *QualificationRun, next RunState) error {
	if r == nil || !r.State.Can(next) {
		return fmt.Errorf("run transition: %w", ErrInvalidState)
	}
	r.State = next
	r.Version++
	return nil
}
func TransitionHold(h *QualityHold, next HoldState, user string) error {
	if h == nil || !h.State.Can(next) {
		return ErrInvalidState
	}
	h.State = next
	h.ClearedBy = &user
	now := timeNow()
	h.ClearedAt = &now
	return nil
}

var timeNow = func() time.Time { return time.Now().UTC() }
