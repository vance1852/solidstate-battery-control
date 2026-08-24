package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"solidstate-battery-control/internal/domain"
	"time"
)

type Manager struct {
	Pool *pgxpool.Pool
	TTL  time.Duration
}

func (m Manager) Create(ctx context.Context, userID string) (string, time.Time, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", time.Time{}, err
	}
	token := hex.EncodeToString(raw)
	exp := time.Now().UTC().Add(m.TTL)
	_, err := m.Pool.Exec(ctx, "INSERT INTO sessions(token_hash,user_id,expires_at) VALUES($1,$2,$3)", hash(token), userID, exp)
	return token, exp, err
}
func hash(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }
func (m Manager) Revoke(ctx context.Context, token string) error {
	tag, err := m.Pool.Exec(ctx, "UPDATE sessions SET revoked_at=now() WHERE token_hash=$1 AND revoked_at IS NULL", hash(token))
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrUnauthorized
	}
	return nil
}
func (m Manager) User(ctx context.Context, token string) (domain.User, error) {
	var u domain.User
	var exp time.Time
	var revoked *time.Time
	err := m.Pool.QueryRow(ctx, "SELECT u.id,u.email,u.role,u.active,s.expires_at,s.revoked_at FROM sessions s JOIN users u ON u.id=s.user_id WHERE s.token_hash=$1", hash(token)).Scan(&u.ID, &u.Email, &u.Role, &u.Active, &exp, &revoked)
	if err != nil {
		return u, fmt.Errorf("session lookup: %w", domain.ErrUnauthorized)
	}
	if !u.Active || revoked != nil || !time.Now().UTC().Before(exp) {
		return u, domain.ErrUnauthorized
	}
	return u, nil
}
func NewUserID() string { return uuid.NewString() }
