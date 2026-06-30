package postgres

import (
	"context"

	service "github.com/example/subscriptions-service/internal/service/subscription"
)

func (r *SubscriptionRepository) TotalCost(ctx context.Context, f service.CostFilter) (int, error) {
	query, args := baseCostArgs(f)

	var total int
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, err
	}
	
	return total, nil
}
