package config

import "os"

type Config struct {
	DatabaseURL, ListenAddr, SessionSecret string
	MigrationDir                           string
}

func Load() Config {
	u := os.Getenv("DATABASE_URL")
	if u == "" {
		u = "postgres://battery:battery@localhost:5432/solidstate_battery?sslmode=disable"
	}
	a := os.Getenv("LISTEN_ADDR")
	if a == "" {
		a = ":8080"
	}
	d := os.Getenv("MIGRATION_DIR")
	if d == "" {
		d = "migrations"
	}
	s := os.Getenv("SESSION_SECRET")
	if s == "" {
		s = "development-secret-change-me"
	}
	return Config{u, a, s, d}
}
