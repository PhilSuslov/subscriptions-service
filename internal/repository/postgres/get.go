package postgres

import (
	"context"
	"errors"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	query := `
	SELECT id, service_name, price, user_id, start_month, end_month, created_at, updated_at
	FROM subscriptions
	WHERE id = $1`

	var subscription domain.Subscription
	row := r.pool.QueryRow(ctx, query, id)
	err := row.Scan(&subscription.ID, &subscription.ServiceName, &subscription.Price, &subscription.UserID,
		&subscription.StartMonth, &subscription.EndMonth, &subscription.CreatedAt, &subscription.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return &subscription, err
}
