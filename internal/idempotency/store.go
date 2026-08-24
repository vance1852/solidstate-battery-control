package idempotency

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Store struct{ Pool *pgxpool.Pool }

func (s Store) Get(ctx context.Context, key string) ([]byte, bool, error) {
	var b []byte
	err := s.Pool.QueryRow(ctx, "SELECT response FROM idempotency_keys WHERE key=$1 AND expires_at>now()", key).Scan(&b)
	if err != nil {
		return nil, false, nil
	}
	return b, true, nil
}
func (s Store) Put(ctx context.Context, key, user string, response any, ttl time.Duration) error {
	b, err := json.Marshal(response)
	if err != nil {
		return err
	}
	_, err = s.Pool.Exec(ctx, "INSERT INTO idempotency_keys(key,actor_id,response,expires_at) VALUES($1,$2,$3,now()+$4::interval) ON CONFLICT(key) DO NOTHING", key, user, b, ttl.String())
	return err
}
