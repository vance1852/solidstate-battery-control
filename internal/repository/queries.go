package repository

import (
	"context"
	"solidstate-battery-control/internal/domain"
)

func (s Store) ListLots(ctx context.Context, states []domain.LotState, limit, offset int) ([]domain.CellLot, error) {
	rows, err := s.Pool.Query(ctx, "SELECT id,lot_code,formulation_id,state,nominal_capacity,version,created_at,updated_at FROM cell_lots ORDER BY updated_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]domain.CellLot, 0)
	for rows.Next() {
		var l domain.CellLot
		if err := rows.Scan(&l.ID, &l.Code, &l.FormulationID, &l.State, &l.Capacity, &l.Version, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}
func (s Store) CountLots(ctx context.Context) (int, error) {
	var n int
	err := s.Pool.QueryRow(ctx, "SELECT count(*) FROM cell_lots").Scan(&n)
	return n, err
}
func (s Store) ListRuns(ctx context.Context, lot string, limit, offset int) ([]domain.QualificationRun, error) {
	rows, err := s.Pool.Query(ctx, "SELECT id,lot_id,COALESCE(module_id::text,''),state,scheduled_at,started_at,finished_at,attempts,version,created_by FROM qualification_runs WHERE lot_id=$1 ORDER BY scheduled_at DESC LIMIT $2 OFFSET $3", lot, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]domain.QualificationRun, 0)
	for rows.Next() {
		var r domain.QualificationRun
		if err := rows.Scan(&r.ID, &r.LotID, &r.ModuleID, &r.State, &r.ScheduledAt, &r.StartedAt, &r.FinishedAt, &r.Attempts, &r.Version, &r.CreatedBy); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
func (s Store) FailRun(ctx context.Context, id string, version int64, reason string) error {
	tag, err := s.Pool.Exec(ctx, "UPDATE qualification_runs SET state='failed',finished_at=now(),version=version+1 WHERE id=$1 AND version=$2", id, version)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (s Store) SucceedRun(ctx context.Context, id string, version int64) error {
	tag, err := s.Pool.Exec(ctx, "UPDATE qualification_runs SET state='succeeded',finished_at=now(),version=version+1 WHERE id=$1 AND version=$2", id, version)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (s Store) TouchRun(ctx context.Context, id string) error {
	_, err := s.Pool.Exec(ctx, "UPDATE qualification_runs SET version=version WHERE id=$1", id)
	return err
}
