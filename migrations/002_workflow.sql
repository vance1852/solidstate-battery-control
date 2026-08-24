CREATE TABLE IF NOT EXISTS audit_events (id uuid PRIMARY KEY, actor_id uuid REFERENCES users(id), entity_type text NOT NULL, entity_id uuid NOT NULL, action text NOT NULL, payload jsonb NOT NULL DEFAULT '{}'::jsonb, request_id text NOT NULL, created_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE IF NOT EXISTS idempotency_keys (key text PRIMARY KEY, actor_id uuid NOT NULL REFERENCES users(id), response jsonb NOT NULL, created_at timestamptz NOT NULL DEFAULT now(), expires_at timestamptz NOT NULL);
CREATE TABLE IF NOT EXISTS worker_jobs (id uuid PRIMARY KEY, kind text NOT NULL, entity_id uuid NOT NULL, state text NOT NULL CHECK(state IN ('pending','running','succeeded','failed')), attempts integer NOT NULL DEFAULT 0, available_at timestamptz NOT NULL DEFAULT now(), last_error text, updated_at timestamptz NOT NULL DEFAULT now(), UNIQUE(kind,entity_id));
CREATE INDEX IF NOT EXISTS idx_runs_due ON qualification_runs(state,scheduled_at);
CREATE INDEX IF NOT EXISTS idx_measurements_run ON measurements(run_id,recorded_at);
CREATE INDEX IF NOT EXISTS idx_audit_entity ON audit_events(entity_type,entity_id,created_at);
CREATE INDEX IF NOT EXISTS idx_jobs_available ON worker_jobs(state,available_at);
