package service

import (
	"context"
	"solidstate-battery-control/internal/domain"
	"solidstate-battery-control/internal/repository"
	"testing"
	"time"
)

func TestWorkflowStatePaths(t *testing.T) {
	w := NewWorkflow(repository.Store{})
	run := domain.QualificationRun{ID: "run-1234", LotID: "lot-1234", State: domain.RunPlanned, ScheduledAt: time.Now()}
	r, err := w.Start(context.Background(), run)
	if err != nil || r.State != domain.RunRunning {
		t.Fatal(err, r)
	}
	r, err = w.Pause(context.Background(), r)
	if err != nil || r.State != domain.RunPaused {
		t.Fatal(err, r)
	}
	r, err = w.Resume(context.Background(), r)
	if err != nil || r.State != domain.RunRunning {
		t.Fatal(err, r)
	}
	r, err = w.Complete(context.Background(), r, true)
	if err != nil || r.State != domain.RunSucceeded || r.FinishedAt == nil {
		t.Fatal(err, r)
	}
}
func TestWorkflowValidation(t *testing.T) {
	w := NewWorkflow(repository.Store{})
	if err := w.ValidateMeasurements(nil); err == nil {
		t.Fatal("empty accepted")
	}
	ms := []domain.Measurement{{Kind: "capacity", Unit: "Ah", Value: 1}}
	if err := w.ValidateMeasurements(ms); err != nil {
		t.Fatal(err)
	}
	ok, _ := w.Qualify(ms)
	if !ok {
		t.Fatal("qualified")
	}
	if err := w.Clear(context.Background(), "", ""); err == nil {
		t.Fatal("empty clear")
	}
}
