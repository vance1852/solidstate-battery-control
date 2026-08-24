package httpapi

import (
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	a := &API{}
	r := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	a.health(w, r)
	if w.Code != 200 {
		t.Fatalf("status %d", w.Code)
	}
}
