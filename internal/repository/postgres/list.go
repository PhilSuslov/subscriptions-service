package postgres

import (
	"context"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
	service "github.com/example/subscriptions-service/internal/service/subscription"
)

func (r *SubscriptionRepository) List(ctx context.Context, f service.ListFilter) ([]domain.Subscription, error) {
	query, args := baseListArgs(f)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Subscription, 0)
	for rows.Next() {
		s, err := scanSubscription(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *s)
	}
	return items, rows.Err()
}
