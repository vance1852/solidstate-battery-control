package service

import (
	"context"
	"solidstate-battery-control/internal/domain"
	"testing"
)

func TestServiceExtended01(t *testing.T) {
	u := domain.User{Active: true, Role: "admin"}
	if Authorize(u, "admin") != nil {
		t.Fatal()
	}
}
func TestServiceExtended02(t *testing.T) {
	u := domain.User{Active: true, Role: "reviewer"}
	if Authorize(u, "reviewer") != nil {
		t.Fatal()
	}
}
func TestServiceExtended03(t *testing.T) {
	u := domain.User{Active: true, Role: "operator"}
	if Authorize(u, "operator") != nil {
		t.Fatal()
	}
}
func TestServiceExtended04(t *testing.T) {
	u := domain.User{Active: true, Role: "operator"}
	if Authorize(u, "admin") == nil {
		t.Fatal()
	}
}
func TestServiceExtended05(t *testing.T) {
	u := domain.User{Active: false, Role: "operator"}
	if Authorize(u, "operator") == nil {
		t.Fatal()
	}
}
func TestServiceExtended06(t *testing.T) {
	ctx := context.Background()
	if RequireOperator(ctx, domain.User{Active: true, Role: "admin"}) != nil {
		t.Fatal()
	}
}
func TestServiceExtended07(t *testing.T) {
	ctx := context.Background()
	if RequireOperator(ctx, domain.User{Active: true, Role: "reviewer"}) == nil {
		t.Fatal()
	}
}
func TestServiceExtended08(t *testing.T) {
	p := Policy{}
	if p.CanCreateLot(domain.User{Active: true, Role: "operator"}) != nil {
		t.Fatal()
	}
}
func TestServiceExtended09(t *testing.T) {
	p := Policy{}
	if p.CanSchedule(domain.User{Active: true, Role: "operator"}) != nil {
		t.Fatal()
	}
}
func TestServiceExtended10(t *testing.T) {
	p := Policy{}
	if p.CanMeasure(domain.User{Active: true, Role: "operator"}) != nil {
		t.Fatal()
	}
}
func TestServiceExtended11(t *testing.T) {
	p := Policy{}
	if p.CanHold(domain.User{Active: true, Role: "reviewer"}) != nil {
		t.Fatal()
	}
}
func TestServiceExtended12(t *testing.T) {
	p := Policy{}
	if p.CanRelease(domain.User{Active: true, Role: "reviewer"}) != nil {
		t.Fatal()
	}
}
func TestServiceExtended13(t *testing.T) {
	p := Policy{}
	if p.CanManageUsers(domain.User{Active: true, Role: "admin"}) != nil {
		t.Fatal()
	}
}
func TestServiceExtended14(t *testing.T) {
	p := Policy{}
	if p.CanHold(domain.User{Active: true, Role: "operator"}) == nil {
		t.Fatal()
	}
}
func TestServiceExtended15(t *testing.T) {
	p := Policy{}
	if p.CanRelease(domain.User{Active: true, Role: "operator"}) == nil {
		t.Fatal()
	}
}
func TestServiceExtended16(t *testing.T) {
	p := Policy{}
	if p.CanManageUsers(domain.User{Active: true, Role: "reviewer"}) == nil {
		t.Fatal()
	}
}
