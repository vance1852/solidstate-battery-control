package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"solidstate-battery-control/internal/domain"
	"time"
)

func (s Store) FindUser(ctx context.Context, email string) (domain.User, string, error) {
	var u domain.User
	var hash string
	err := s.Pool.QueryRow(ctx, "SELECT id,email,password_hash,role,active FROM users WHERE email=$1", email).Scan(&u.ID, &u.Email, &hash, &u.Role, &u.Active)
	return u, hash, err
}
func (s Store) CreateUser(ctx context.Context, u domain.User, hash string) error {
	_, err := s.Pool.Exec(ctx, "INSERT INTO users(id,email,password_hash,role,active) VALUES($1,$2,$3,$4,$5)", u.ID, u.Email, hash, u.Role, u.Active)
	return err
}
func (s Store) CreateSession(ctx context.Context, token, user string, exp time.Time) error {
	_, err := s.Pool.Exec(ctx, "INSERT INTO sessions(token_hash,user_id,expires_at) VALUES($1,$2,$3)", token, user, exp)
	return err
}
func (s Store) RevokeSession(ctx context.Context, token string) error {
	tag, err := s.Pool.Exec(ctx, "UPDATE sessions SET revoked_at=now() WHERE token_hash=$1", token)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
func (s Store) EnsureSeed(ctx context.Context, hash string) error {
	_, err := s.Pool.Exec(ctx, "INSERT INTO users(id,email,password_hash,role) VALUES('00000000-0000-0000-0000-000000000001','admin@battery.local',$1,'admin') ON CONFLICT(email) DO NOTHING", hash)
	return err
}

var _ = pgx.ErrNoRows
