package subscription

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseMonth(t *testing.T) {
	m, err := ParseMonth("07-2025")
	require.NoError(t, err)
	require.Equal(t, "07-2025", m.String())
}

func TestParseMonthInvalidFormat(t *testing.T) {
	_, err := ParseMonth("2025-07")
	require.Error(t, err)
}
