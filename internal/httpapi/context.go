package httpapi

import (
	"context"
	"time"
)

type deadlineKey struct{}

func withRequestTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, d)
}
func requestDeadline(ctx context.Context) time.Time {
	if d, ok := ctx.Deadline(); ok {
		return d
	}
	return time.Time{}
}
func canceled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
func valueOr(ctx context.Context, key any, fallback string) string {
	if v, ok := ctx.Value(key).(string); ok {
		return v
	}
	return fallback
}
func contextRequestID(ctx context.Context) string { return valueOr(ctx, requestKey{}, "") }
