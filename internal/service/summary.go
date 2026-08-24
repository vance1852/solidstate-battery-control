package service

import (
	"context"
	"solidstate-battery-control/internal/domain"
	"solidstate-battery-control/internal/repository"
)

type Summary struct {
	Lot          domain.CellLot
	Runs         []domain.QualificationRun
	Measurements []domain.Measurement
	Audits       []domain.AuditEvent
}

func (s Service) BuildSummary(ctx context.Context, lotID string) (Summary, error) {
	lot, err := s.Store.GetLot(ctx, lotID)
	if err != nil {
		return Summary{}, err
	}
	audits, err := s.Store.ListAudit(ctx, "cell_lot", lotID, 50, 0)
	if err != nil {
		return Summary{}, err
	}
	return Summary{Lot: lot, Audits: audits}, nil
}

var _ = repository.RunFilter{}
