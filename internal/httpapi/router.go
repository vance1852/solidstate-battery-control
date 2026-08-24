package httpapi

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"solidstate-battery-control/internal/domain"
	"solidstate-battery-control/internal/service"
	"strings"
	"time"
)

type API struct {
	Svc  service.Service
	Auth service.AuthService
	Pool *pgxpool.Pool
}

func New(s service.Service, a service.AuthService, p *pgxpool.Pool) *API {
	return &API{Svc: s, Auth: a, Pool: p}
}
func (a *API) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", a.health)
	mux.HandleFunc("GET /readyz", a.ready)
	mux.HandleFunc("POST /api/login", a.login)
	mux.HandleFunc("POST /api/logout", a.logout)
	mux.HandleFunc("POST /api/lots", a.createLot)
	mux.HandleFunc("GET /api/lots/", a.getLot)
	mux.HandleFunc("POST /api/runs", a.createRun)
	mux.HandleFunc("POST /api/measurements", a.measure)
	mux.HandleFunc("POST /api/holds", a.openHold)
	mux.HandleFunc("POST /api/holds/clear", a.clearHold)
	return requestMiddleware(mux)
}
func requestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}
		ctx := context.WithValue(r.Context(), requestKey{}, id)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type requestKey struct{}

func requestID(ctx context.Context) string { v, _ := ctx.Value(requestKey{}).(string); return v }
func (a *API) health(w http.ResponseWriter, _ *http.Request) {
	write(w, 200, map[string]any{"status": "ok"})
}
func (a *API) ready(w http.ResponseWriter, r *http.Request) {
	if err := a.Pool.Ping(r.Context()); err != nil {
		writeError(w, 503, err)
		return
	}
	write(w, 200, map[string]any{"status": "ready"})
}
func (a *API) login(w http.ResponseWriter, r *http.Request) {
	var in struct{ Email, Password string }
	if !decode(w, r, &in) {
		return
	}
	token, user, exp, err := a.Auth.Login(r.Context(), in.Email, in.Password)
	if err != nil {
		writeError(w, 401, err)
		return
	}
	write(w, 200, map[string]any{"token": token, "expires_at": exp, "user": user})
}
func (a *API) logout(w http.ResponseWriter, r *http.Request) {
	if err := a.Auth.Logout(r.Context(), bearer(r)); err != nil {
		writeError(w, 401, err)
		return
	}
	write(w, 204, nil)
}
func (a *API) user(r *http.Request) (domain.User, error) {
	token := bearer(r)
	if token == "" {
		return domain.User{}, domain.ErrUnauthorized
	}
	return a.AuthUser(r.Context(), token)
}
func (a *API) AuthUser(ctx context.Context, token string) (domain.User, error) {
	var u domain.User
	err := a.Pool.QueryRow(ctx, "SELECT u.id,u.email,u.role,u.active FROM sessions s JOIN users u ON u.id=s.user_id WHERE s.token_hash=$1 AND s.revoked_at IS NULL AND s.expires_at>now()", hashToken(token)).Scan(&u.ID, &u.Email, &u.Role, &u.Active)
	if err != nil {
		return u, domain.ErrUnauthorized
	}
	return u, nil
}
func (a *API) createLot(w http.ResponseWriter, r *http.Request) {
	u, err := a.user(r)
	if err != nil || service.Authorize(u, "operator", "admin") != nil {
		writeError(w, 401, domain.ErrUnauthorized)
		return
	}
	var in struct {
		Code, FormulationID string
		Capacity            float64
	}
	if !decode(w, r, &in) {
		return
	}
	lot := domain.CellLot{ID: uuid.NewString(), Code: in.Code, FormulationID: in.FormulationID, State: domain.LotDraft, Capacity: in.Capacity}
	if err := a.Svc.CreateLot(r.Context(), lot); err != nil {
		writeError(w, 400, err)
		return
	}
	write(w, 201, lot)
}
func (a *API) getLot(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/lots/")
	lot, err := a.Svc.GetLot(r.Context(), id)
	if errors.Is(err, domain.ErrNotFound) {
		writeError(w, 404, err)
		return
	}
	if err != nil {
		writeError(w, 500, err)
		return
	}
	write(w, 200, lot)
}
func (a *API) createRun(w http.ResponseWriter, r *http.Request) {
	u, err := a.user(r)
	if err != nil {
		writeError(w, 401, err)
		return
	}
	if err = service.RequireOperator(r.Context(), u); err != nil {
		writeError(w, 403, err)
		return
	}
	var in struct {
		LotID, ModuleID string
		ScheduledAt     time.Time
	}
	if !decode(w, r, &in) {
		return
	}
	run := domain.QualificationRun{ID: uuid.NewString(), LotID: in.LotID, ModuleID: in.ModuleID, State: domain.RunPlanned, ScheduledAt: in.ScheduledAt, CreatedBy: u.ID}
	if err := a.Svc.CreateRun(r.Context(), run); err != nil {
		writeError(w, 400, err)
		return
	}
	write(w, 201, run)
}
func (a *API) measure(w http.ResponseWriter, r *http.Request) {
	u, err := a.user(r)
	if err != nil {
		writeError(w, 401, err)
		return
	}
	var in struct {
		RunID, Kind, Unit string
		Value             float64
	}
	if !decode(w, r, &in) {
		return
	}
	m := domain.Measurement{ID: uuid.NewString(), RunID: in.RunID, Kind: in.Kind, Unit: in.Unit, Value: in.Value, RecordedBy: u.ID}
	if err := a.Svc.AddMeasurement(r.Context(), m); err != nil {
		writeError(w, 400, err)
		return
	}
	write(w, 201, m)
}
func (a *API) openHold(w http.ResponseWriter, r *http.Request) {
	u, err := a.user(r)
	if err != nil {
		writeError(w, 401, err)
		return
	}
	if err = service.Authorize(u, "reviewer", "admin"); err != nil {
		writeError(w, 403, err)
		return
	}
	var in struct{ LotID, Reason string }
	if !decode(w, r, &in) {
		return
	}
	h := domain.QualityHold{ID: uuid.NewString(), LotID: in.LotID, Reason: in.Reason, State: domain.HoldOpen, OpenedBy: u.ID}
	if err := a.Svc.OpenHold(r.Context(), h, requestID(r.Context())); err != nil {
		writeError(w, 400, err)
		return
	}
	write(w, 201, h)
}
func (a *API) clearHold(w http.ResponseWriter, r *http.Request) {
	u, err := a.user(r)
	if err != nil {
		writeError(w, 401, err)
		return
	}
	if err = service.Authorize(u, "reviewer", "admin"); err != nil {
		writeError(w, 403, err)
		return
	}
	var in struct{ ID string }
	if !decode(w, r, &in) {
		return
	}
	if err := a.Svc.ClearHold(r.Context(), in.ID, u.ID); err != nil {
		writeError(w, 400, err)
		return
	}
	write(w, 204, nil)
}
func bearer(r *http.Request) string {
	return strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
}
func hashToken(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }
func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(v); err != nil {
		writeError(w, 400, err)
		return false
	}
	return true
}
func write(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	if v != nil {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(v)
	}
}
func writeError(w http.ResponseWriter, status int, err error) {
	code := "internal_error"
	if errors.Is(err, domain.ErrUnauthorized) {
		code = "unauthorized"
	}
	if errors.Is(err, domain.ErrForbidden) {
		code = "forbidden"
	}
	if errors.Is(err, domain.ErrValidation) {
		code = "invalid"
	}
	if errors.Is(err, domain.ErrNotFound) {
		code = "not_found"
	}
	write(w, status, map[string]any{"error": code, "message": err.Error()})
}
