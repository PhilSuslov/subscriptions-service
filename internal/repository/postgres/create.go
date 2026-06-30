package postgres

import (
	"context"
	"time"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
)

func (r *SubscriptionRepository) Create(ctx context.Context, s *domain.Subscription) error {
	query := `
	INSERT INTO subscriptions (id, service_name, price, user_id, start_month, end_month)
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING created_at, updated_at`

	var endMonth *time.Time
	if s.EndMonth != nil {
		endMonth = &s.EndMonth.Time
	}
	return r.pool.QueryRow(ctx, query, s.ID, s.ServiceName, s.Price, s.UserID, s.StartMonth.Time, endMonth).Scan(&s.CreatedAt, &s.UpdatedAt)
}
