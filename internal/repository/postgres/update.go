package postgres

import (
	"context"
	"errors"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
	"github.com/jackc/pgx/v5"
)

func (r *SubscriptionRepository) Update(ctx context.Context, s *domain.Subscription) error {
	query := `UPDATE subscriptions
	SET service_name=$2, price=$3, user_id=$4, start_month=$5, end_month=$6, updated_at=now()
	WHERE id=$1
	RETURNING created_at, updated_at`

	var end any
	if s.EndMonth != nil {
		end = s.EndMonth.Time
	}
	err := r.pool.QueryRow(ctx, query, s.ID, s.ServiceName, s.Price, s.UserID, s.StartMonth.Time, end).Scan(&s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}
