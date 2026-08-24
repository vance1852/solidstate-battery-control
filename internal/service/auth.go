package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"solidstate-battery-control/internal/domain"
	"solidstate-battery-control/internal/repository"
	"strings"
	"time"
)

type AuthService struct {
	Store repository.Store
	TTL   time.Duration
}

func (a AuthService) Login(ctx context.Context, email, password string) (string, domain.User, time.Time, error) {
	u, stored, err := a.Store.FindUser(ctx, strings.TrimSpace(strings.ToLower(email)))
	if err != nil || !u.Active || bcrypt.CompareHashAndPassword([]byte(stored), []byte(password)) != nil {
		return "", u, time.Time{}, domain.ErrUnauthorized
	}
	token, exp, err := newToken()
	if err != nil {
		return "", u, time.Time{}, err
	}
	if err = a.Store.CreateSession(ctx, hashToken(token), u.ID, exp); err != nil {
		return "", u, time.Time{}, fmt.Errorf("create session: %w", err)
	}
	return token, u, exp, nil
}
func (a AuthService) Logout(ctx context.Context, token string) error {
	return a.Store.RevokeSession(ctx, hashToken(token))
}
func newToken() (string, time.Time, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", time.Time{}, err
	}
	return hex.EncodeToString(b), time.Now().UTC().Add(24 * time.Hour), nil
}
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
