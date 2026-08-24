package service

import (
	"context"
	"fmt"
	"solidstate-battery-control/internal/domain"
)

type Policy struct{}

func (Policy) CanCreateLot(u domain.User) error   { return Authorize(u, "operator", "admin") }
func (Policy) CanSchedule(u domain.User) error    { return Authorize(u, "operator", "admin") }
func (Policy) CanMeasure(u domain.User) error     { return Authorize(u, "operator", "admin") }
func (Policy) CanHold(u domain.User) error        { return Authorize(u, "reviewer", "admin") }
func (Policy) CanRelease(u domain.User) error     { return Authorize(u, "reviewer", "admin") }
func (Policy) CanManageUsers(u domain.User) error { return Authorize(u, "admin") }
func (Policy) CheckContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
func (p Policy) CheckLot(ctx context.Context, u domain.User, l domain.CellLot) error {
	if err := p.CheckContext(ctx); err != nil {
		return err
	}
	if err := p.CanCreateLot(u); err != nil {
		return err
	}
	if l.Code == "" {
		return fmt.Errorf("code missing: %w", domain.ErrValidation)
	}
	return nil
}
func (p Policy) CheckRun(ctx context.Context, u domain.User, r domain.QualificationRun) error {
	if err := p.CheckContext(ctx); err != nil {
		return err
	}
	if err := p.CanSchedule(u); err != nil {
		return err
	}
	if r.LotID == "" {
		return domain.ErrValidation
	}
	return nil
}
func (p Policy) CheckMeasurement(ctx context.Context, u domain.User, m domain.Measurement) error {
	if err := p.CheckContext(ctx); err != nil {
		return err
	}
	if err := p.CanMeasure(u); err != nil {
		return err
	}
	return domain.ValidateMeasurement(m.Kind, m.Unit, m.Value)
}
func (p Policy) CheckHold(ctx context.Context, u domain.User, h domain.QualityHold) error {
	if err := p.CheckContext(ctx); err != nil {
		return err
	}
	if err := p.CanHold(u); err != nil {
		return err
	}
	if h.Reason == "" {
		return domain.ErrValidation
	}
	return nil
}
