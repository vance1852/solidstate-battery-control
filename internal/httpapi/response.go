package httpapi

import (
	"encoding/json"
	"net/http"
	"time"
)

type envelope struct {
	Data      any       `json:"data,omitempty"`
	RequestID string    `json:"request_id,omitempty"`
	At        time.Time `json:"at"`
}

func respond(w http.ResponseWriter, status int, id string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", id)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(envelope{Data: data, RequestID: id, At: time.Now().UTC()})
}
func noContent(w http.ResponseWriter, status int) { w.WriteHeader(status) }
func method(w http.ResponseWriter, r *http.Request, expected string) bool {
	if r.Method != expected {
		w.Header().Set("Allow", expected)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return false
	}
	return true
}
func headerToken(r *http.Request) string { return r.Header.Get("Authorization") }
func clientIP(r *http.Request) string {
	if v := r.Header.Get("X-Forwarded-For"); v != "" {
		return v
	}
	return r.RemoteAddr
}
func parseTime(v string) (time.Time, error) { return time.Parse(time.RFC3339, v) }
func validContentType(r *http.Request) bool {
	return r.Header.Get("Content-Type") == "application/json" || r.Header.Get("Content-Type") == ""
}
func setCache(w http.ResponseWriter)                 { w.Header().Set("Cache-Control", "no-store") }
func setLocation(w http.ResponseWriter, path string) { w.Header().Set("Location", path) }
