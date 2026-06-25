package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/example/subscriptions-service/internal/config"
	"github.com/stretchr/testify/require"
)

func TestLoadAppliesEnvOverrides(t *testing.T) {
	t.Setenv("HTTP_ADDR", ":9090")
	t.Setenv("POSTGRES_DSN", "postgres://user:pass@localhost:5432/db?sslmode=disable")
	t.Setenv("POSTGRES_MAX_CONNS", "20")
	t.Setenv("POSTGRES_MIN_CONNS", "2")
	t.Setenv("POSTGRES_MAX_CONN_LIFETIME", "30m")

	cfg, err := config.Load("")
	require.NoError(t, err)
	require.Equal(t, ":9090", cfg.HTTP.Addr)
	require.Equal(t, "postgres://user:pass@localhost:5432/db?sslmode=disable", cfg.Postgres.DSN)
	require.Equal(t, int32(20), cfg.Postgres.MaxConns)
	require.Equal(t, int32(2), cfg.Postgres.MinConns)
	require.Equal(t, 30*time.Minute, cfg.Postgres.MaxConnLifetime)
}

func TestLoadRequiresDSN(t *testing.T) {
	_ = os.Unsetenv("POSTGRES_DSN")
	_, err := config.Load("")
	require.Error(t, err)
}
