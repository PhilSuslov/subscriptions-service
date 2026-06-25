package postgres

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBuildCommonFilters(t *testing.T) {
	userID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")
	serviceName := " Yandex Plus "

	where, args := buildCommonFilters(&userID, &serviceName)

	require.Equal(t, "WHERE user_id=$1 AND service_name ILIKE $2", where)
	require.Len(t, args, 2)
	require.Equal(t, userID, args[0])
	require.Equal(t, "Yandex Plus", args[1])
}

func TestBuildCommonFiltersEmpty(t *testing.T) {
	where, args := buildCommonFilters(nil, nil)

	require.Empty(t, where)
	require.Empty(t, args)
}
