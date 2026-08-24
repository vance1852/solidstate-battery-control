package repository

import (
	"context"
	"github.com/google/uuid"
	"solidstate-battery-control/internal/domain"
)

func (s Store) CreateModule(ctx context.Context, m domain.Module) error {
	_, err := s.Pool.Exec(ctx, "INSERT INTO modules(id,serial,lot_id,state) VALUES($1,$2,$3,$4)", m.ID, m.Serial, m.LotID, m.State)
	return err
}
func (s Store) GetModule(ctx context.Context, id string) (domain.Module, error) {
	var m domain.Module
	err := s.Pool.QueryRow(ctx, "SELECT id,serial,lot_id,state,version FROM modules WHERE id=$1", id).Scan(&m.ID, &m.Serial, &m.LotID, &m.State, &m.Version)
	return m, err
}
func (s Store) ChangeModuleState(ctx context.Context, id string, version int64, state string) error {
	tag, err := s.Pool.Exec(ctx, "UPDATE modules SET state=$1,version=version+1 WHERE id=$2 AND version=$3", state, id, version)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return domain.ErrConflict
	}
	return nil
}
func NewModule(serial, lot string) domain.Module {
	return domain.Module{ID: uuid.NewString(), Serial: serial, LotID: lot, State: "assembled"}
}
func (s Store) ListModules(ctx context.Context, lot string) ([]domain.Module, error) {
	rows, err := s.Pool.Query(ctx, "SELECT id,serial,lot_id,state,version FROM modules WHERE lot_id=$1 ORDER BY serial", lot)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]domain.Module, 0)
	for rows.Next() {
		var m domain.Module
		if err := rows.Scan(&m.ID, &m.Serial, &m.LotID, &m.State, &m.Version); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
func (s Store) CountModules(ctx context.Context, lot string) (int, error) {
	var n int
	err := s.Pool.QueryRow(ctx, "SELECT count(*) FROM modules WHERE lot_id=$1", lot).Scan(&n)
	return n, err
}
