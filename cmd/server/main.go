package main

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"solidstate-battery-control/internal/config"
	"solidstate-battery-control/internal/db"
	"solidstate-battery-control/internal/httpapi"
	"solidstate-battery-control/internal/repository"
	"solidstate-battery-control/internal/service"
	"solidstate-battery-control/internal/worker"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	pool, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	if err = db.Migrate(ctx, pool, cfg.MigrationDir); err != nil {
		log.Fatal(err)
	}
	seed, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	store := repository.Store{Pool: pool}
	if err = store.EnsureSeed(ctx, string(seed)); err != nil {
		log.Fatal(err)
	}
	svc := service.New(store)
	api := httpapi.New(svc, service.AuthService{Store: store, TTL: 24 * time.Hour}, pool)
	w := worker.New(svc, slog.Default())
	w.Start(ctx)
	srv := &http.Server{Addr: cfg.ListenAddr, Handler: api.Handler(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	go func() {
		<-ctx.Done()
		shutdown, stop := context.WithTimeout(context.Background(), 10*time.Second)
		defer stop()
		w.Stop()
		_ = srv.Shutdown(shutdown)
	}()
	log.Printf("solid-state battery control listening on %s", cfg.ListenAddr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
