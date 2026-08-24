package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"solidstate-battery-control/internal/domain"
)

func (s Store) CreateLot(ctx context.Context, lot domain.CellLot) error {
	_, err := s.Pool.Exec(ctx, "INSERT INTO cell_lots(id,lot_code,formulation_id,state,nominal_capacity) VALUES($1,$2,$3,$4,$5)", lot.ID, lot.Code, lot.FormulationID, lot.State, lot.Capacity)
	return err
}
func (s Store) GetLot(ctx context.Context, id string) (domain.CellLot, error) {
	var l domain.CellLot
	err := s.Pool.QueryRow(ctx, "SELECT id,lot_code,formulation_id,state,nominal_capacity,version,created_at,updated_at FROM cell_lots WHERE id=$1", id).Scan(&l.ID, &l.Code, &l.FormulationID, &l.State, &l.Capacity, &l.Version, &l.CreatedAt, &l.UpdatedAt)
	if err == pgx.ErrNoRows {
		return l, domain.ErrNotFound
	}
	return l, err
}
func (s Store) TransitionLot(ctx context.Context, id string, expected int64, next domain.LotState) (domain.CellLot, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return domain.CellLot{}, err
	}
	defer tx.Rollback(ctx)
	var l domain.CellLot
	err = tx.QueryRow(ctx, "SELECT id,lot_code,formulation_id,state,nominal_capacity,version,created_at,updated_at FROM cell_lots WHERE id=$1 FOR UPDATE", id).Scan(&l.ID, &l.Code, &l.FormulationID, &l.State, &l.Capacity, &l.Version, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return l, err
	}
	if l.Version != expected {
		return l, domain.ErrConflict
	}
	if err = domain.TransitionLot(&l, next); err != nil {
		return l, err
	}
	tag, err := tx.Exec(ctx, "UPDATE cell_lots SET state=$1,version=$2,updated_at=now() WHERE id=$3 AND version=$4", l.State, l.Version, id, expected)
	if err != nil || tag.RowsAffected() != 1 {
		return l, domain.ErrConflict
	}
	if err = tx.Commit(ctx); err != nil {
		return l, err
	}
	return l, nil
}
func NewLot(code, form string, cap float64) domain.CellLot {
	return domain.CellLot{ID: uuid.NewString(), Code: code, FormulationID: form, State: domain.LotDraft, Capacity: cap}
}
func requireLot(l domain.CellLot) error {
	if l.ID == "" {
		return fmt.Errorf("id required")
	}
	return nil
}
