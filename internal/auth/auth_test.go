package auth

import (
	"strings"
	"testing"
)

func TestPasswordRoundTrip(t *testing.T) {
	h, err := HashPassword("secret")
	if err != nil {
		t.Fatal(err)
	}
	if !CheckPassword(h, "secret") || CheckPassword(h, "wrong") {
		t.Fatal("password verification")
	}
	if !strings.HasPrefix(h, "$2") {
		t.Fatal("bcrypt hash")
	}
}
