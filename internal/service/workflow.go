package service

import (
	"context"
	"fmt"
	"solidstate-battery-control/internal/domain"
	"solidstate-battery-control/internal/repository"
	"time"
)

type Workflow struct{ Store repository.Store }

func NewWorkflow(s repository.Store) Workflow { return Workflow{Store: s} }
func (w Workflow) Schedule(ctx context.Context, lot domain.CellLot, user string, when time.Time) (domain.QualificationRun, error) {
	if err := domain.EnsureID(lot.ID); err != nil {
		return domain.QualificationRun{}, err
	}
	if !lot.CanSchedule() {
		return domain.QualificationRun{}, fmt.Errorf("lot cannot schedule: %w", domain.ErrInvalidState)
	}
	run := repository.NewRun(lot.ID, user, when)
	if err := w.Store.CreateRun(ctx, run); err != nil {
		return run, fmt.Errorf("schedule run: %w", err)
	}
	return run, nil
}
func (w Workflow) Start(ctx context.Context, run domain.QualificationRun) (domain.QualificationRun, error) {
	if err := domain.TransitionRun(&run, domain.RunRunning); err != nil {
		return run, err
	}
	return run, nil
}
func (w Workflow) Pause(ctx context.Context, run domain.QualificationRun) (domain.QualificationRun, error) {
	if err := domain.TransitionRun(&run, domain.RunPaused); err != nil {
		return run, err
	}
	return run, nil
}
func (w Workflow) Resume(ctx context.Context, run domain.QualificationRun) (domain.QualificationRun, error) {
	if err := domain.TransitionRun(&run, domain.RunRunning); err != nil {
		return run, err
	}
	return run, nil
}
func (w Workflow) Complete(ctx context.Context, run domain.QualificationRun, ok bool) (domain.QualificationRun, error) {
	next := domain.RunFailed
	if ok {
		next = domain.RunSucceeded
	}
	if err := domain.TransitionRun(&run, next); err != nil {
		return run, err
	}
	now := time.Now().UTC()
	run.FinishedAt = &now
	return run, nil
}
func (w Workflow) ValidateMeasurements(ms []domain.Measurement) error {
	if len(ms) == 0 {
		return fmt.Errorf("measurements required: %w", domain.ErrValidation)
	}
	for _, m := range ms {
		if err := domain.ValidateMeasurement(m.Kind, m.Unit, m.Value); err != nil {
			return err
		}
	}
	return nil
}
func (w Workflow) Qualify(ms []domain.Measurement) (bool, []string) {
	rules := domain.DefaultRules()
	return rules.Evaluate(ms)
}
func (w Workflow) Release(ctx context.Context, lot domain.CellLot, version int64) (domain.CellLot, error) {
	if lot.Version != version {
		return lot, domain.ErrConflict
	}
	return w.Store.TransitionLot(ctx, lot.ID, version, domain.LotReleased)
}
func (w Workflow) Hold(ctx context.Context, lotID, reason, user, request string) error {
	h := repository.NewHold(lotID, reason, user)
	return w.Store.OpenHold(ctx, h, request)
}
func (w Workflow) Clear(ctx context.Context, holdID, user string) error {
	if holdID == "" || user == "" {
		return domain.ErrValidation
	}
	return w.Store.ClearHold(ctx, holdID, user)
}
func (w Workflow) Due(ctx context.Context) ([]domain.QualificationRun, error) {
	return w.Store.StartDue(ctx, time.Now().UTC())
}
