# Solid-State Battery Pilot Control

Production backend for coordinating solid-state battery cell lots, module assembly pilots, qualification runs, thermal measurements, quality holds, and release decisions. The service uses Go 1.22, PostgreSQL 16, versioned SQL migrations, transactional workflows, optimistic concurrency, idempotency keys, audit events, and a restart-safe worker.

Run locally with `docker compose up -d postgres`, `go test ./... -count=1`, and `go run ./cmd/server`. The default database URL is `postgres://battery:battery@localhost:5432/solidstate_battery?sslmode=disable`. Health endpoints are `/healthz` and `/readyz`; authentication is session based with login, expiry and revocation.
