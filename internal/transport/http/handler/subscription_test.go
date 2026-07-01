package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestParseListFilter(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/subscriptions?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service_name=Yandex+Plus&limit=12&offset=3", nil)

	f, err := parseListFilter(req)
	require.NoError(t, err)
	require.NotNil(t, f.UserID)
	require.Equal(t, 12, f.Limit)
	require.Equal(t, 3, f.Offset)
	require.NotNil(t, f.ServiceName)
}

func TestParseListFilterInvalidUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/subscriptions?user_id=bad", nil)
	_, err := parseListFilter(req)
	require.Error(t, err)
}

func TestParseListFilterDefaults(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/subscriptions", nil)

	f, err := parseListFilter(req)
	require.NoError(t, err)
	require.Equal(t, 50, f.Limit)
	require.Equal(t, 0, f.Offset)
	require.Nil(t, f.UserID)
	require.Nil(t, f.ServiceName)
}

func TestParseListFilterBlankServiceName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/subscriptions?service_name=++", nil)

	f, err := parseListFilter(req)
	require.NoError(t, err)
	require.Nil(t, f.ServiceName)
}

func TestToUpdateCommand(t *testing.T) {
	id := uuid.New()
	req := subscriptionRequest{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		StartDate:   "07-2025",
	}
	cmd, err := toUpdateCommand(id, req)
	require.NoError(t, err)
	require.Equal(t, id, cmd.ID)
}

func TestToUpdateCommandBadUUID(t *testing.T) {
	_, err := toUpdateCommand(uuid.New(), subscriptionRequest{UserID: "bad"})
	require.Error(t, err)
}
