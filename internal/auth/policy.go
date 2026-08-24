package auth

import (
	"solidstate-battery-control/internal/domain"
	"strings"
	"time"
)

func NormalizeEmail(v string) string { return strings.ToLower(strings.TrimSpace(v)) }
func NormalizeRole(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	if domain.ValidateRole(v) {
		return v
	}
	return "operator"
}
func SessionValid(exp time.Time, revoked *time.Time, now time.Time) bool {
	if revoked != nil {
		return false
	}
	return now.Before(exp)
}
func Can(role, action string) bool {
	switch action {
	case "lot:create", "run:create", "measurement:add":
		return role == "operator" || role == "admin"
	case "hold:open", "hold:clear", "lot:release":
		return role == "reviewer" || role == "admin"
	case "user:manage":
		return role == "admin"
	}
	return false
}
func Actions(role string) []string {
	out := make([]string, 0)
	for _, a := range []string{"lot:create", "run:create", "measurement:add", "hold:open", "hold:clear", "lot:release", "user:manage"} {
		if Can(role, a) {
			out = append(out, a)
		}
	}
	return out
}
