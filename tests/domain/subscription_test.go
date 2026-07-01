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

func TestValidateRejectsEmptyServiceName(t *testing.T) {
	start, _ := domain.ParseMonth("07-2025")
	s, err := domain.New("   ", 400, uuid.New(), start, nil)
	require.ErrorIs(t, err, domain.ErrInvalidServiceName)
	require.Nil(t, s)
}

func TestValidateRejectsNilUser(t *testing.T) {
	start, _ := domain.ParseMonth("07-2025")
	_, err := domain.New("Yandex Plus", 400, uuid.Nil, start, nil)
	require.ErrorIs(t, err, domain.ErrInvalidUserID)
}

func TestValidateRejectsZeroPrice(t *testing.T) {
	start, _ := domain.ParseMonth("07-2025")
	s := &domain.Subscription{ServiceName: "Yandex Plus", Price: 0, UserID: uuid.New(), StartMonth: start}
	require.ErrorIs(t, s.Validate(), domain.ErrInvalidPrice)
}

func TestValidateRejectsEndBeforeStart(t *testing.T) {
	start, _ := domain.ParseMonth("07-2025")
	end, _ := domain.ParseMonth("06-2025")
	s := &domain.Subscription{ServiceName: "Yandex Plus", Price: 400, UserID: uuid.New(), StartMonth: start, EndMonth: &end}
	require.ErrorIs(t, s.Validate(), domain.ErrInvalidPeriod)
}
