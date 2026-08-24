package repository

import (
	"context"
	"github.com/google/uuid"
	"solidstate-battery-control/internal/domain"
)

func (s Store) OpenHold(ctx context.Context, h domain.QualityHold, requestID string) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, "UPDATE cell_lots SET state='hold',version=version+1,updated_at=now() WHERE id=$1", h.LotID); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, "INSERT INTO quality_holds(id,lot_id,reason,state,opened_by) VALUES($1,$2,$3,'open',$4)", h.ID, h.LotID, h.Reason, h.OpenedBy); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, "INSERT INTO audit_events(id,actor_id,entity_type,entity_id,action,request_id) VALUES($1,$2,'cell_lot',$3,'hold_opened',$4)", uuid.NewString(), h.OpenedBy, h.LotID, requestID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (s Store) ClearHold(ctx context.Context, id, user string) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var lot string
	if err = tx.QueryRow(ctx, "SELECT lot_id FROM quality_holds WHERE id=$1 AND state='open' FOR UPDATE", id).Scan(&lot); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, "UPDATE quality_holds SET state='cleared',cleared_by=$1,cleared_at=now() WHERE id=$2", user, id); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, "UPDATE cell_lots SET state='qualified',version=version+1,updated_at=now() WHERE id=$1", lot); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func NewHold(lot, reason, user string) domain.QualityHold {
	return domain.QualityHold{ID: uuid.NewString(), LotID: lot, Reason: reason, OpenedBy: user, State: domain.HoldOpen}
}
