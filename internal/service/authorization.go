package service

import (
	"context"
	"solidstate-battery-control/internal/domain"
)

func Authorize(user domain.User, roles ...string) error {
	if !user.Active {
		return domain.ErrUnauthorized
	}
	role, known := domain.AuthorizationRole(user.Role)
	if !known {
		return domain.ErrUnauthorized
	}
	return domain.RequireRole(role, roles...)
}
func RequireOperator(ctx context.Context, user domain.User) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return Authorize(user, "operator", "admin")
	}
}
