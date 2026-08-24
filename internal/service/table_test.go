package service

import (
	"context"
	"solidstate-battery-control/internal/domain"
	"solidstate-battery-control/internal/repository"
	"testing"
	"time"
)

func TestServiceTable01(t *testing.T) {
	if New(repository.Store{}).CreateLot(context.Background(), domain.CellLot{Code: "", Capacity: 1}) == nil {
		t.Fatal()
	}
}
func TestServiceTable02(t *testing.T) {
	if New(repository.Store{}).CreateLot(context.Background(), domain.CellLot{Code: "L", Capacity: 0}) == nil {
		t.Fatal()
	}
}
func TestServiceTable03(t *testing.T) {
	if New(repository.Store{}).AddMeasurement(context.Background(), domain.Measurement{Kind: "x", Unit: "u", Value: 1}) == nil {
		t.Fatal()
	}
}
func TestServiceTable04(t *testing.T) {
	if New(repository.Store{}).AddMeasurement(context.Background(), domain.Measurement{Kind: "capacity", Unit: "", Value: 1}) == nil {
		t.Fatal()
	}
}
func TestServiceTable05(t *testing.T) {
	if New(repository.Store{}).AddMeasurement(context.Background(), domain.Measurement{Kind: "capacity", Unit: "Ah", Value: -1}) == nil {
		t.Fatal()
	}
}
func TestServiceTable06(t *testing.T) {
	if Authorize(domain.User{Active: true, Role: "admin"}, "admin") != nil {
		t.Fatal()
	}
}
func TestServiceTable07(t *testing.T) {
	if Authorize(domain.User{Active: true, Role: "operator"}, "admin") == nil {
		t.Fatal()
	}
}
func TestServiceTable08(t *testing.T) {
	if Authorize(domain.User{Active: false, Role: "admin"}, "admin") == nil {
		t.Fatal()
	}
}
func TestServiceTable09(t *testing.T) {
	if RequireOperator(context.Background(), domain.User{Active: true, Role: "operator"}) != nil {
		t.Fatal()
	}
}
func TestServiceTable10(t *testing.T) {
	if RequireOperator(context.Background(), domain.User{Active: true, Role: "reviewer"}) == nil {
		t.Fatal()
	}
}
func TestServiceTable11(t *testing.T) {
	ctx, c := context.WithCancel(context.Background())
	c()
	if RequireOperator(ctx, domain.User{Active: true, Role: "operator"}) == nil {
		t.Fatal()
	}
}
func TestServiceTable12(t *testing.T) {
	w := NewWorkflow(repository.Store{})
	if _, err := w.Start(context.Background(), domain.QualificationRun{State: domain.RunSucceeded}); err == nil {
		t.Fatal()
	}
}
func TestServiceTable13(t *testing.T) {
	w := NewWorkflow(repository.Store{})
	if _, err := w.Pause(context.Background(), domain.QualificationRun{State: domain.RunPlanned}); err == nil {
		t.Fatal()
	}
}
func TestServiceTable14(t *testing.T) {
	w := NewWorkflow(repository.Store{})
	if _, err := w.Resume(context.Background(), domain.QualificationRun{State: domain.RunRunning}); err == nil {
		t.Fatal()
	}
}
func TestServiceTable15(t *testing.T) {
	w := NewWorkflow(repository.Store{})
	if _, err := w.Complete(context.Background(), domain.QualificationRun{State: domain.RunRunning}, true); err != nil {
		t.Fatal()
	}
}
func TestServiceTable16(t *testing.T) {
	w := NewWorkflow(repository.Store{})
	if _, err := w.Complete(context.Background(), domain.QualificationRun{State: domain.RunPlanned}, false); err != nil {
		t.Fatal()
	}
}
func TestServiceTable17(t *testing.T) {
	r := repository.NewRun("lot-1234", "user-1234", time.Now())
	if r.State != domain.RunPlanned {
		t.Fatal()
	}
}
func TestServiceTable18(t *testing.T) {
	m := repository.NewMeasurement("run-1234", "capacity", "Ah", 1, "user-1234")
	if m.ID == "" {
		t.Fatal()
	}
}
func TestServiceTable19(t *testing.T) {
	h := repository.NewHold("lot-1234", "drift", "user-1234")
	if h.State != domain.HoldOpen {
		t.Fatal()
	}
}
func TestServiceTable20(t *testing.T) {
	f := repository.NewFormulation("oxide", "sulfide", 1)
	if f.ID == "" || f.Approved {
		t.Fatal()
	}
}
func TestServiceTable21(t *testing.T) {
	m := repository.NewModule("M-1", "lot-1234")
	if m.State != "assembled" {
		t.Fatal()
	}
}
func TestServiceTable22(t *testing.T) {
	if IsNotFound(domain.ErrNotFound) == false {
		t.Fatal()
	}
}
func TestServiceTable23(t *testing.T) {
	if IsNotFound(domain.ErrConflict) {
		t.Fatal()
	}
}
func TestServiceTable24(t *testing.T) {
	if New(repository.Store{}).CreateRun(context.Background(), domain.QualificationRun{}) == nil {
		t.Fatal()
	}
}
func TestServiceTable25(t *testing.T) {
	if New(repository.Store{}).OpenHold(context.Background(), domain.QualityHold{}, "") == nil {
		t.Fatal()
	}
}
