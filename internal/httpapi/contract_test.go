package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareRequestID(t *testing.T) {
	h := requestMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requestID(r.Context()) == "" {
			t.Fatal("missing id")
		}
		w.WriteHeader(204)
	}))
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Header().Get("X-Request-ID") == "" {
		t.Fatal("response id")
	}
}
func TestBearer(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer abc")
	if bearer(r) != "abc" {
		t.Fatal("token")
	}
}
