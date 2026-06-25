package subscription

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewSubscriptionSuccess(t *testing.T) {
	start, err := ParseMonth("07-2025")
	require.NoError(t, err)
	s, err := New("Yandex Plus", 400, uuid.New(), start, nil)
	require.NoError(t, err)
	require.Equal(t, "Yandex Plus", s.ServiceName)
	require.Equal(t, 400, s.Price)
}

func TestNewSubscriptionInvalidPrice(t *testing.T) {
	start, _ := ParseMonth("07-2025")
	_, err := New("Yandex Plus", 0, uuid.New(), start, nil)
	require.ErrorIs(t, err, ErrInvalidPrice)
}

func TestNewSubscriptionInvalidPeriod(t *testing.T) {
	start, _ := ParseMonth("07-2025")
	end, _ := ParseMonth("06-2025")
	_, err := New("Yandex Plus", 400, uuid.New(), start, &end)
	require.ErrorIs(t, err, ErrInvalidPeriod)
}
