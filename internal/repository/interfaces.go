package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"solidstate-battery-control/internal/domain"
	"time"
)

type Tx interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Commit(context.Context) error
	Rollback(context.Context) error
}
type Store struct{ Pool *pgxpool.Pool }
type TxStore struct{ Tx pgx.Tx }
type RunFilter struct {
	State         string
	Limit, Offset int
}
type UserStore interface {
	FindUser(context.Context, string) (domain.User, string, error)
	CreateSession(context.Context, string, string, time.Time) error
	RevokeSession(context.Context, string) error
}
