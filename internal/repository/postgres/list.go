package postgres

import (
	"context"
	"fmt"
	"strings"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
	service "github.com/example/subscriptions-service/internal/service/subscription"
)

func (r *SubscriptionRepository) List(ctx context.Context, f service.ListFilter) ([]domain.Subscription, error) {

	args, clauses := buildBaseFilters(f.UserID, f.ServiceName)
	args = append(args, f.Limit, f.Offset)

	query := `
	SELECT id, service_name, price, user_id, start_month, end_month, created_at, updated_at
	FROM subscriptions`

	if len(clauses) > 0 {
		query += "\nWHERE " + strings.Join(clauses, " AND ")
	}
	query += fmt.Sprintf("\nORDER BY created_at DESC, id DESC\nLIMIT $%d OFFSET $%d", len(args)-1, len(args))

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
