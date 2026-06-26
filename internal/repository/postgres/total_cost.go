package postgres

import (
	"context"
	"fmt"
	"strings"

	service "github.com/example/subscriptions-service/internal/service/subscription"
)

func (r *SubscriptionRepository) TotalCost(ctx context.Context, f service.CostFilter) (int, error) {
	args, clauses := buildBaseFilters(f.UserID, f.ServiceName)
	args = append(args, f.PeriodFrom.Time, f.PeriodTo.Time)

	periodFromArg := len(args) - 1
	periodToArg := len(args)
	clauses = append(clauses, fmt.Sprintf("start_month <= $%d", periodToArg), fmt.Sprintf("(end_month IS NULL OR end_month >= $%d)", periodFromArg))

	where := "WHERE " + strings.Join(clauses, " AND ")
	query := fmt.Sprintf(`
	SELECT COALESCE(SUM(
		price * (
			((EXTRACT(YEAR FROM LEAST(COALESCE(end_month, $%d), $%d))::int - EXTRACT(YEAR FROM GREATEST(start_month, $%d))::int) * 12) +
			(EXTRACT(MONTH FROM LEAST(COALESCE(end_month, $%d), $%d))::int - EXTRACT(MONTH FROM GREATEST(start_month, $%d))::int) + 1
		)
	), 0)::int
	FROM subscriptions
	%s`,
		periodToArg, periodToArg, periodFromArg, periodToArg, periodToArg, periodFromArg, where,
	)

	var total int
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}
