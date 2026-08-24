package repository

import (
	"context"
	"github.com/google/uuid"
	"solidstate-battery-control/internal/domain"
)

func (s Store) CreateFormulation(ctx context.Context, f domain.Formulation) error {
	_, err := s.Pool.Exec(ctx, "INSERT INTO electrolyte_formulations(id,name,version,chemistry,approved) VALUES($1,$2,$3,$4,$5)", f.ID, f.Name, f.Version, f.Chemistry, f.Approved)
	return err
}
func (s Store) ApproveFormulation(ctx context.Context, id string) error {
	tag, err := s.Pool.Exec(ctx, "UPDATE electrolyte_formulations SET approved=true WHERE id=$1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return domain.ErrNotFound
	}
	return nil
}
func (s Store) GetFormulation(ctx context.Context, id string) (domain.Formulation, error) {
	var f domain.Formulation
	err := s.Pool.QueryRow(ctx, "SELECT id,name,version,chemistry,approved FROM electrolyte_formulations WHERE id=$1", id).Scan(&f.ID, &f.Name, &f.Version, &f.Chemistry, &f.Approved)
	return f, err
}
func NewFormulation(name, chem string, version int) domain.Formulation {
	return domain.Formulation{ID: uuid.NewString(), Name: name, Chemistry: chem, Version: version}
}
func (s Store) ListFormulations(ctx context.Context, approved bool) ([]domain.Formulation, error) {
	rows, err := s.Pool.Query(ctx, "SELECT id,name,version,chemistry,approved FROM electrolyte_formulations WHERE approved=$1 ORDER BY name,version", approved)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]domain.Formulation, 0)
	for rows.Next() {
		var f domain.Formulation
		if err := rows.Scan(&f.ID, &f.Name, &f.Version, &f.Chemistry, &f.Approved); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}
func (s Store) DeleteFormulation(ctx context.Context, id string) error {
	tag, err := s.Pool.Exec(ctx, "DELETE FROM electrolyte_formulations WHERE id=$1 AND approved=false", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return domain.ErrConflict
	}
	return nil
}
