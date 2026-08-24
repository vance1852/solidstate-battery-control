package service

import (
	"solidstate-battery-control/internal/domain"
	"testing"
)

func TestInvalidCompletionDoesNotStampFinishedAt(t *testing.T) {
	run := domain.QualificationRun{ID: "run-010", State: domain.RunSucceeded}
	workflow := Workflow{}
	updated, err := workflow.Complete(nil, run, true)
	if err == nil {
		t.Fatal("completion of an already finished run unexpectedly succeeded")
	}
	if updated.FinishedAt != nil {
		t.Fatal("failed completion polluted the run with a finished timestamp")
	}
}
