package repository

import (
	"context"
	"github.com/google/uuid"
	"solidstate-battery-control/internal/domain"
)

func (s Store) AddMeasurement(ctx context.Context, m domain.Measurement) error {
	_, err := s.Pool.Exec(ctx, "INSERT INTO measurements(id,run_id,kind,value,unit,recorded_by) VALUES($1,$2,$3,$4,$5,$6)", m.ID, m.RunID, m.Kind, m.Value, m.Unit, m.RecordedBy)
	return err
}
func (s Store) ListMeasurements(ctx context.Context, run string) ([]domain.Measurement, error) {
	rows, err := s.Pool.Query(ctx, "SELECT id,run_id,kind,value,unit,recorded_at,recorded_by FROM measurements WHERE run_id=$1 ORDER BY recorded_at", run)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Measurement
	for rows.Next() {
		var m domain.Measurement
		if err := rows.Scan(&m.ID, &m.RunID, &m.Kind, &m.Value, &m.Unit, &m.RecordedAt, &m.RecordedBy); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
func NewMeasurement(run, kind, unit string, value float64, user string) domain.Measurement {
	return domain.Measurement{ID: uuid.NewString(), RunID: run, Kind: kind, Unit: unit, Value: value, RecordedBy: user}
}
