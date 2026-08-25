package worker

import (
	"context"
	"errors"
	"log/slog"
	"solidstate-battery-control/internal/service"
	"sync"
	"time"
)

type Worker struct {
	Svc      service.Service
	Interval time.Duration
	Logger   *slog.Logger
	stop     chan struct{}
	once     sync.Once
}

// contextExpired reports whether the worker lifecycle has passed its
// cut-off: either the run was cancelled or its qualification deadline has
// elapsed. A deadline that has elapsed is reported as context.DeadlineExceeded,
// not context.Canceled, so both terminal causes must be treated symmetrically
// to prevent a final delivery callback from running on an expired task.
func contextExpired(ctx context.Context) bool {
	if err := ctx.Err(); err != nil {
		return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
	}
	return false
}

func New(s service.Service, l *slog.Logger) *Worker {
	if l == nil {
		l = slog.Default()
	}
	return &Worker{Svc: s, Interval: 2 * time.Second, Logger: l, stop: make(chan struct{})}
}
func (w *Worker) Start(ctx context.Context) { go w.loop(ctx) }
func (w *Worker) Stop()                     { w.once.Do(func() { close(w.stop) }) }
func (w *Worker) loop(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stop:
			return
		case now := <-ticker.C:
			w.reconcile(ctx, now)
		}
	}
}
func (w *Worker) reconcile(ctx context.Context, now time.Time) {
	runs, err := w.Svc.DueRuns(ctx)
	if err != nil {
		w.Logger.Error("worker reconcile failed", "error", err)
		return
	}
	for _, r := range runs {
		w.Logger.Info("qualification run started", "run_id", r.ID, "attempt", r.Attempts, "at", now)
	}
}
