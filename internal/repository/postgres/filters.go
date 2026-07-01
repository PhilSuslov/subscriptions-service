package postgres

import (
	"strconv"
	"strings"

	service "github.com/example/subscriptions-service/internal/service/subscription"
	"github.com/google/uuid"
)

func normalizedServiceName(serviceName *string) (string, bool) {
	if serviceName == nil {
		return "", false
	}
	v := strings.TrimSpace(*serviceName)
	return v, v != ""
}

func placeholder(n int) string {
	return "$" + strconv.Itoa(n)
}

func buildWhereClause(userID *uuid.UUID, serviceName *string) (string, []any) {
	clauses := make([]string, 0, 2)
	args := make([]any, 0, 2)

	if userID != nil {
		clauses = append(clauses, "user_id="+placeholder(len(args)+1))
		args = append(args, *userID)
	}

	if name, ok := normalizedServiceName(serviceName); ok {
		clauses = append(clauses, "service_name ILIKE "+placeholder(len(args)+1))
		args = append(args, name)
	}

	if len(clauses) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(clauses, " AND "), args
}

func buildListQuery(f service.ListFilter) (string, []any) {
	where, args := buildWhereClause(f.UserID, f.ServiceName)

	base := `
	SELECT id, service_name, price, user_id, start_month, end_month, created_at, updated_at
	FROM subscriptions`

	if where != "" {
		base += "\n" + where
	}

	switch len(args) {
	case 0:
		base += "\nORDER BY created_at DESC, id DESC\nLIMIT $1 OFFSET $2"
	case 1:
		base += "\nORDER BY created_at DESC, id DESC\nLIMIT $2 OFFSET $3"
	case 2:
		base += "\nORDER BY created_at DESC, id DESC\nLIMIT $3 OFFSET $4"
	}

	args = append(args, f.Limit, f.Offset)
	return base, args
}

func buildCostQuery(f service.CostFilter) (string, []any) {
	where, args := buildWhereClause(f.UserID, f.ServiceName)
	base := `
	SELECT COALESCE(SUM(
		price * (
			((EXTRACT(YEAR FROM LEAST(COALESCE(end_month, $END), $END))::int - EXTRACT(YEAR FROM GREATEST(start_month, $FROM))::int) * 12) +
			(EXTRACT(MONTH FROM LEAST(COALESCE(end_month, $END), $END))::int - EXTRACT(MONTH FROM GREATEST(start_month, $FROM))::int + 1)
		)
	), 0)::int
	FROM subscriptions`

	switch len(args) {
	case 0:
		base = strings.ReplaceAll(base, "$FROM", "1")
		base = strings.ReplaceAll(base, "$END", "2")
	case 1:
		base = strings.ReplaceAll(base, "$FROM", "2")
		base = strings.ReplaceAll(base, "$END", "3")
	case 2:
		base = strings.ReplaceAll(base, "$FROM", "3")
		base = strings.ReplaceAll(base, "$END", "4")
	}

	if where != "" {
		base += "\n" + where + " AND start_month <= " + placeholder(len(args)+1) +
			" AND (end_month IS NULL OR end_month >= " + placeholder(len(args)+2) + ")"
	} else {
		base += "\nWHERE start_month <= " + placeholder(len(args)+1) +
			" AND (end_month IS NULL OR end_month >= " + placeholder(len(args)+2) + ")"
	}

	args = append(args, f.PeriodFrom.Time, f.PeriodTo.Time)
	return base, args
}

func BuildCommonFiltersForTest(userID *uuid.UUID, serviceName *string) (string, []any) {
	return buildWhereClause(userID, serviceName)
}

func baseListArgs(f service.ListFilter) (string, []any) {
	return buildListQuery(f)
}

func baseCostArgs(f service.CostFilter) (string, []any) {
	return buildCostQuery(f)
}
