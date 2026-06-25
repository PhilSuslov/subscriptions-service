package subscription_test

import (
	"testing"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewSubscriptionSuccess(t *testing.T) {
	start, err := domain.ParseMonth("07-2025")
	require.NoError(t, err)
	s, err := domain.New("Yandex Plus", 400, uuid.New(), start, nil)
	require.NoError(t, err)
	require.Equal(t, "Yandex Plus", s.ServiceName)
	require.Equal(t, 400, s.Price)
}

func TestNewSubscriptionInvalidPrice(t *testing.T) {
	start, _ := domain.ParseMonth("07-2025")
	_, err := domain.New("Yandex Plus", 0, uuid.New(), start, nil)
	require.ErrorIs(t, err, domain.ErrInvalidPrice)
}

func TestNewSubscriptionInvalidPeriod(t *testing.T) {
	start, _ := domain.ParseMonth("07-2025")
	end, _ := domain.ParseMonth("06-2025")
	_, err := domain.New("Yandex Plus", 400, uuid.New(), start, &end)
	require.ErrorIs(t, err, domain.ErrInvalidPeriod)
}
