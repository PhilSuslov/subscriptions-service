package subscription_test

import (
	"testing"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
	"github.com/stretchr/testify/require"
)

func TestParseMonth(t *testing.T) {
	m, err := domain.ParseMonth("07-2025")
	require.NoError(t, err)
	require.Equal(t, "07-2025", m.String())
}

func TestParseMonthInvalidFormat(t *testing.T) {
	_, err := domain.ParseMonth("2025-07")
	require.Error(t, err)
}
