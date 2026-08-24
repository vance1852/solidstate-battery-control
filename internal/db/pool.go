package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func Open(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 12
	cfg.MinConns = 2
	cfg.MaxConnLifetime = 30 * time.Minute
	return pgxpool.NewWithConfig(ctx, cfg)
}
