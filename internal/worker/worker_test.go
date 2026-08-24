package worker

import (
	"context"
	"solidstate-battery-control/internal/service"
	"testing"
	"time"
)

func TestStopIsIdempotent(t *testing.T) { w := New(service.Service{}, nil); w.Stop(); w.Stop() }
func TestContextDone(t *testing.T) {
	ctx, c := context.WithCancel(context.Background())
	c()
	w := New(service.Service{}, nil)
	w.Start(ctx)
	time.Sleep(time.Millisecond * 5)
	w.Stop()
}
