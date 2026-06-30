package postgres

import (
	"strconv"
	"strings"
	"time"

	service "github.com/example/subscriptions-service/internal/service/subscription"
	"github.com/google/uuid"
)

type subscriptionQueryFilter struct {
	clauses []string
	args    []any
}

func newSubscriptionQueryFilter(userID *uuid.UUID, serviceName *string) subscriptionQueryFilter {
	filter := subscriptionQueryFilter{
		clauses: make([]string, 0, 2),
		args:    make([]any, 0, 2),
	}

	if userID != nil {
		filter.clauses = append(filter.clauses, "user_id="+placeholder(len(filter.args)+1))
		filter.args = append(filter.args, *userID)
	}

	if name, ok := normalizedServiceName(serviceName); ok {
		filter.clauses = append(filter.clauses, "service_name ILIKE "+placeholder(len(filter.args)+1))
		filter.args = append(filter.args, name)
	}

	return filter
}

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

func (f subscriptionQueryFilter) where() string {
	if len(f.clauses) == 0 {
		return ""
	}
	return "WHERE " + strings.Join(f.clauses, " AND ")
}

func (f subscriptionQueryFilter) listQuery() (string, []interface{}) {
	query := `
SELECT id, service_name, price, user_id, start_month, end_month, created_at, updated_at
FROM subscriptions`
	if where := f.where(); where != "" {
		query += "\n" + where
	}
	query += "\nORDER BY created_at DESC, id DESC\nLIMIT " + placeholder(len(f.args)+1) + " OFFSET " + placeholder(len(f.args)+2)
	return query, append([]any{}, f.args...)
}

func (f subscriptionQueryFilter) costQuery(periodFrom, periodTo time.Time) (string, []interface{}) {
	base := len(f.args)
	query := `
	SELECT COALESCE(SUM(
		price * (
			((EXTRACT(YEAR FROM LEAST(COALESCE(end_month, ` + placeholder(base+2) + `), ` + placeholder(base+2) + `))::int - EXTRACT(YEAR FROM GREATEST(start_month, ` + placeholder(base+1) + `))::int) * 12) +
			(EXTRACT(MONTH FROM LEAST(COALESCE(end_month, ` + placeholder(base+2) + `), ` + placeholder(base+2) + `))::int - EXTRACT(MONTH FROM GREATEST(start_month, ` + placeholder(base+1) + `))::int + 1)
		)
	), 0)::int
	FROM subscriptions`
	if where := f.where(); where != "" {
		query += "\n" + where + " AND start_month <= " + placeholder(base+1) + " AND (end_month IS NULL OR end_month >= " + placeholder(base+2) + ")"
	} else {
		query += "\nWHERE start_month <= " + placeholder(base+1) + " AND (end_month IS NULL OR end_month >= " + placeholder(base+2) + ")"
	}
	args := append(append([]any{}, f.args...), periodFrom, periodTo)
	return query, args
}

func BuildCommonFiltersForTest(userID *uuid.UUID, serviceName *string) (string, []any) {
	filter := newSubscriptionQueryFilter(userID, serviceName)
	return filter.where(), filter.args
}

func baseListArgs(f service.ListFilter) (string, []any) {
	filter := newSubscriptionQueryFilter(f.UserID, f.ServiceName)
	query, args := filter.listQuery()
	args = append(args, f.Limit, f.Offset)
	return query, args
}

func baseCostArgs(f service.CostFilter) (string, []any) {
	filter := newSubscriptionQueryFilter(f.UserID, f.ServiceName)
	query, args := filter.costQuery(f.PeriodFrom.Time, f.PeriodTo.Time)
	return query, args
}
