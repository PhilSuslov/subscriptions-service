package postgres

import (
	"context"
	"testing"

	"github.com/example/subscriptions-service/internal/config"
	"github.com/stretchr/testify/require"
)

func TestNewPoolRejectsInvalidDSN(t *testing.T) {
	_, err := NewPool(context.Background(), config.PostgresConfig{DSN: "://bad"})
	require.Error(t, err)
}
