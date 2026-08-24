package audit

import (
	"solidstate-battery-control/internal/repository"
	"testing"
)

func TestPayload(t *testing.T) {
	if len(repository.AuditPayload(map[string]string{"a": "b"})) == 0 {
		t.Fatal("empty payload")
	}
}
