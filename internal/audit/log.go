package audit

import (
	"context"
	"log/slog"
	"solidstate-battery-control/internal/domain"
	"solidstate-battery-control/internal/repository"
)

type Logger struct {
	Store repository.Store
	Log   *slog.Logger
}

func (l Logger) Record(ctx context.Context, actor, entity, id, action, request string, payload any) error {
	e := domain.AuditEvent{ActorID: actor, EntityType: entity, EntityID: id, Action: action, RequestID: request, Payload: repository.AuditPayload(payload)}
	if err := l.Store.Audit(ctx, e); err != nil {
		if l.Log != nil {
			l.Log.Error("audit write failed", "error", err)
		}
		return err
	}
	return nil
}
