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
	clauses = append(clauses,
		fmt.Sprintf("start_month <= $%d", periodToArg),
		fmt.Sprintf("(end_month IS NULL OR end_month >= $%d)", periodFromArg),
	)

	query := fmt.Sprintf(`
	SELECT COALESCE(SUM(
		price * (
			(EXTRACT(YEAR FROM age(LEAST(COALESCE(end_month, $%d), $%d), GREATEST(start_month, $%d)))::int * 12) +
			EXTRACT(MONTH FROM age(LEAST(COALESCE(end_month, $%d), $%d), GREATEST(start_month, $%d)))::int + 1
		)
	), 0)::int
	FROM subscriptions
	WHERE %s`,
		periodToArg,
		periodToArg,
		periodFromArg,
		periodToArg,
		periodToArg,
		periodFromArg,
		strings.Join(clauses, " AND "),
	)

	var total int
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}
