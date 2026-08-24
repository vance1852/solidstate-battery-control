package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Options struct {
	MaxBodyBytes    int64
	ShutdownTimeout time.Duration
	WorkerInterval  time.Duration
	LogJSON         bool
}

func DefaultOptions() Options {
	return Options{MaxBodyBytes: 1 << 20, ShutdownTimeout: 10 * time.Second, WorkerInterval: 2 * time.Second}
}
func LoadOptions() Options {
	o := DefaultOptions()
	if v, err := strconv.ParseInt(os.Getenv("MAX_BODY_BYTES"), 10, 64); err == nil && v > 0 {
		o.MaxBodyBytes = v
	}
	if v, err := time.ParseDuration(os.Getenv("WORKER_INTERVAL")); err == nil && v > 0 {
		o.WorkerInterval = v
	}
	o.LogJSON = os.Getenv("LOG_FORMAT") == "json"
	return o
}
func (o Options) Validate() error {
	if o.MaxBodyBytes < 1024 {
		return fmt.Errorf("body limit too small")
	}
	if o.ShutdownTimeout <= 0 {
		return fmt.Errorf("shutdown timeout required")
	}
	if o.WorkerInterval <= 0 {
		return fmt.Errorf("worker interval required")
	}
	return nil
}
func EnvBool(name string) bool { return os.Getenv(name) == "1" || os.Getenv(name) == "true" }
