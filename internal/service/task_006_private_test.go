package service

import (
	"testing"
	"time"
)

func TestNewSessionExpiryIsInFuture(t *testing.T) {
	before := time.Now().UTC()
	_, expires, err := newToken()
	if err != nil {
		t.Fatal(err)
	}
	if !expires.After(before) {
		t.Fatalf("new session already expired at %s", expires)
	}
	if expires.Sub(before) < 23*time.Hour {
		t.Fatalf("session lifetime is too short: %s", expires.Sub(before))
	}
}
