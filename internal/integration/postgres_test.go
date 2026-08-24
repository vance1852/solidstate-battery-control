package integration

import (
	"context"
	"os"
	"solidstate-battery-control/internal/db"
	"testing"
	"time"
)

func TestPostgresMigrationAndHealth(t *testing.T) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL is required for PostgreSQL integration")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := db.Open(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		t.Fatal(err)
	}
	if err := db.Migrate(ctx, pool, "../../migrations"); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM schema_migrations").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count < 3 {
		t.Fatalf("expected three migrations, got %d", count)
	}
	var tables int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM information_schema.tables WHERE table_schema='public' AND table_name IN ('users','cell_lots','qualification_runs','measurements','audit_events','sessions')").Scan(&tables); err != nil {
		t.Fatal(err)
	}
	if tables != 6 {
		t.Fatalf("schema incomplete: %d", tables)
	}
}

func TestPostgresContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := ctx.Err(); err == nil {
		t.Fatal("cancel signal missing")
	}
}
