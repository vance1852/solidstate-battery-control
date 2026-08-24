package worker

import (
	"context"
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

func contextExpired(ctx context.Context) bool {
	return ctx.Err() == context.Canceled
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
