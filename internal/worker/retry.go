package worker

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type RetryPolicy struct {
	MaxAttempts int
	BaseDelay   time.Duration
}

func DefaultRetry() RetryPolicy { return RetryPolicy{MaxAttempts: 4, BaseDelay: time.Second} }
func (p RetryPolicy) Delay(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	d := p.BaseDelay
	for i := 1; i < attempt; i++ {
		d *= 2
	}
	return d
}
func retryable(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) {
		return true
	}
	return true
}
func (p RetryPolicy) ShouldRetry(attempt int, err error) bool {
	return retryable(err) && attempt < p.MaxAttempts
}
func (p RetryPolicy) Run(ctx context.Context, fn func(context.Context) error) error {
	var last error
	for attempt := 1; attempt <= p.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := fn(ctx); err == nil {
			return nil
		} else {
			last = err
			if !p.ShouldRetry(attempt, err) {
				break
			}
			timer := time.NewTimer(p.Delay(attempt))
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		}
	}
	return fmt.Errorf("retry exhausted: %w", last)
}
func Backoff(attempt int) time.Duration { return DefaultRetry().Delay(attempt) }
func IsPermanent(err error) bool        { return err != nil && errors.Is(err, context.Canceled) }
