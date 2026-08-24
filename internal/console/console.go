package console

import (
	"context"
	"fmt"
	"solidstate-battery-control/internal/domain"
)

type Console struct{}

func (Console) CanPublish(ctx context.Context, user domain.User, lot domain.CellLot) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if user.Role != "admin" && user.Role != "reviewer" {
			return domain.ErrForbidden
		}
		if lot.State != "qualified" {
			return fmt.Errorf("lot is not qualified: %w", domain.ErrInvalidState)
		}
		return nil
	}
}
