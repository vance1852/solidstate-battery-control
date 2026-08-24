package httpapi

import (
	"net/http/httptest"
	"testing"
)

func TestResponseEnvelope(t *testing.T) {
	w := httptest.NewRecorder()
	respond(w, 200, "req-1", map[string]string{"ok": "yes"})
	if w.Code != 200 || w.Header().Get("X-Request-ID") != "req-1" {
		t.Fatal()
	}
}
func TestMethod(t *testing.T) {
	r := httptest.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	if method(w, r, "GET") {
		t.Fatal()
	}
	if w.Code != 405 || w.Header().Get("Allow") != "GET" {
		t.Fatal()
	}
}
func TestParseTime(t *testing.T) {
	if _, err := parseTime("bad"); err == nil {
		t.Fatal()
	}
	if _, err := parseTime("2026-01-01T00:00:00Z"); err != nil {
		t.Fatal(err)
	}
}
