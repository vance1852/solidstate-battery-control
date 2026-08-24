package service

import (
	"context"
	"errors"
	"fmt"
	"solidstate-battery-control/internal/domain"
	"solidstate-battery-control/internal/repository"
	"time"
)

type Service struct{ Store repository.Store }

func New(s repository.Store) Service { return Service{Store: s} }
func (s Service) CreateLot(ctx context.Context, lot domain.CellLot) error {
	if err := domain.ValidateLot(lot.Code, lot.Capacity); err != nil {
		return err
	}
	if lot.State == "" {
		lot.State = domain.LotDraft
	}
	return s.Store.CreateLot(ctx, lot)
}
func (s Service) GetLot(ctx context.Context, id string) (domain.CellLot, error) {
	return s.Store.GetLot(ctx, id)
}
func (s Service) ChangeLotState(ctx context.Context, id string, version int64, next domain.LotState) (domain.CellLot, error) {
	l, err := s.Store.TransitionLot(ctx, id, version, next)
	if err != nil {
		return l, fmt.Errorf("change lot state: %w", err)
	}
	return l, nil
}
func (s Service) CreateRun(ctx context.Context, r domain.QualificationRun) error {
	if r.LotID == "" || r.CreatedBy == "" {
		return domain.ErrValidation
	}
	return s.Store.CreateRun(ctx, r)
}
func (s Service) AddMeasurement(ctx context.Context, m domain.Measurement) error {
	if err := domain.ValidateMeasurement(m.Kind, m.Unit, m.Value); err != nil {
		return err
	}
	return s.Store.AddMeasurement(ctx, m)
}
func (s Service) OpenHold(ctx context.Context, h domain.QualityHold, requestID string) error {
	if h.LotID == "" || h.Reason == "" {
		return domain.ErrValidation
	}
	return s.Store.OpenHold(ctx, h, requestID)
}
func (s Service) ClearHold(ctx context.Context, id, user string) error {
	return s.Store.ClearHold(ctx, id, user)
}
func (s Service) DueRuns(ctx context.Context) ([]domain.QualificationRun, error) {
	return s.Store.StartDue(ctx, time.Now().UTC())
}
func IsNotFound(err error) bool { return errors.Is(err, domain.ErrNotFound) }
