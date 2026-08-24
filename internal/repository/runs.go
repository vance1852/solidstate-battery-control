package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"solidstate-battery-control/internal/domain"
	"time"
)

func (s Store) CreateRun(ctx context.Context, r domain.QualificationRun) error {
	_, err := s.Pool.Exec(ctx, "INSERT INTO qualification_runs(id,lot_id,module_id,state,scheduled_at,created_by) VALUES($1,$2,NULLIF($3,''),$4,$5,$6)", r.ID, r.LotID, r.ModuleID, r.State, r.ScheduledAt, r.CreatedBy)
	return err
}
func (s Store) GetRun(ctx context.Context, id string) (domain.QualificationRun, error) {
	var r domain.QualificationRun
	err := s.Pool.QueryRow(ctx, "SELECT id,lot_id,COALESCE(module_id::text,''),state,scheduled_at,started_at,finished_at,attempts,version,created_by FROM qualification_runs WHERE id=$1", id).Scan(&r.ID, &r.LotID, &r.ModuleID, &r.State, &r.ScheduledAt, &r.StartedAt, &r.FinishedAt, &r.Attempts, &r.Version, &r.CreatedBy)
	if err == pgx.ErrNoRows {
		return r, domain.ErrNotFound
	}
	return r, err
}
func (s Store) StartDue(ctx context.Context, now time.Time) ([]domain.QualificationRun, error) {
	rows, err := s.Pool.Query(ctx, "UPDATE qualification_runs SET state='running',started_at=COALESCE(started_at,now()),attempts=attempts+1,version=version+1 WHERE state='planned' AND scheduled_at<= $1 RETURNING id,lot_id,COALESCE(module_id::text,''),state,scheduled_at,started_at,finished_at,attempts,version,created_by", now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.QualificationRun
	for rows.Next() {
		var r domain.QualificationRun
		if err := rows.Scan(&r.ID, &r.LotID, &r.ModuleID, &r.State, &r.ScheduledAt, &r.StartedAt, &r.FinishedAt, &r.Attempts, &r.Version, &r.CreatedBy); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
func NewRun(lot, user string, when time.Time) domain.QualificationRun {
	return domain.QualificationRun{ID: uuid.NewString(), LotID: lot, State: domain.RunPlanned, ScheduledAt: when, CreatedBy: user}
}
