package repository

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"solidstate-battery-control/internal/domain"
)

func (s Store) Audit(ctx context.Context, e domain.AuditEvent) error {
	if e.ID == "" {
		e.ID = uuid.NewString()
	}
	if len(e.Payload) == 0 {
		e.Payload = []byte(`{}`)
	}
	_, err := s.Pool.Exec(ctx, "INSERT INTO audit_events(id,actor_id,entity_type,entity_id,action,payload,request_id) VALUES($1,NULLIF($2,'')::uuid,$3,$4,$5,$6,$7)", e.ID, e.ActorID, e.EntityType, e.EntityID, e.Action, e.Payload, e.RequestID)
	return err
}
func (s Store) ListAudit(ctx context.Context, entity, id string, limit, offset int) ([]domain.AuditEvent, error) {
	rows, err := s.Pool.Query(ctx, "SELECT id,COALESCE(actor_id::text,''),entity_type,entity_id,action,payload,request_id,created_at FROM audit_events WHERE entity_type=$1 AND entity_id=$2 ORDER BY created_at DESC LIMIT $3 OFFSET $4", entity, id, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.AuditEvent
	for rows.Next() {
		var e domain.AuditEvent
		if err := rows.Scan(&e.ID, &e.ActorID, &e.EntityType, &e.EntityID, &e.Action, &e.Payload, &e.RequestID, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
func AuditPayload(v any) []byte { b, _ := json.Marshal(v); return b }
